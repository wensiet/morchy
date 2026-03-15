#!/usr/bin/env python3
"""Simplified E2E tests - single file, no dependencies beyond requests/psycopg2."""

import os
import socket
import subprocess
import sys
import time
import traceback
from pathlib import Path

import psycopg2
import requests

os.environ["DOCKER_HOST"] = "unix:///Users/wensiet/.colima/default/docker.sock"

PROJECT_ROOT = Path(__file__).parent.parent.parent
BIN_DIR = PROJECT_ROOT / "bin"
MIGRATIONS_DIR = PROJECT_ROOT / "migrations"

DB_HOST = "localhost"
DB_PORT = 5433
DB_USER = "postgres"
DB_PASSWORD = "postgres"
DB_NAME = "morchy_test"

CP_PORT = 8080
CP_API_URL = f"http://localhost:{CP_PORT}/api/v1"

TIMEOUT = 60
POLL_INTERVAL = 1


def log(msg):
    print(f"[{time.strftime('%H:%M:%S')}] {msg}")


def run_cmd(cmd, check=True, capture_output=True, env=None):
    if capture_output:
        result = subprocess.run(cmd, capture_output=True, text=True, env=env, check=False)
        if check and result.returncode != 0:
            log(f"Command failed: {' '.join(cmd)}")
            log(f"STDERR: {result.stderr}")
            sys.exit(1)
        return result
    subprocess.run(cmd, check=check, env=env)


def wait_for_port(host, port, timeout=30):
    start = time.time()
    while time.time() - start < timeout:
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.settimeout(0.5)
            result = sock.connect_ex((host, port))
            sock.close()
            if result == 0:
                return None
        except (OSError, socket.timeout):
            pass
        time.sleep(0.2)
    raise RuntimeError(f"Port {port} not open after {timeout}s")


def db_conn():
    return psycopg2.connect(
        host=DB_HOST, port=DB_PORT, user=DB_USER, password=DB_PASSWORD, dbname=DB_NAME
    )


def start_postgres():
    log("Starting PostgreSQL...")
    result = run_cmd(["docker", "ps", "-aq", "-f", "name=morchy-e2e-postgres"], check=False)
    if result.stdout.strip():
        log("PostgreSQL container exists, stopping...")
        run_cmd(["docker", "rm", "-f", "morchy-e2e-postgres"], check=False)
        time.sleep(3)

    run_cmd(
        [
            "docker",
            "run",
            "-d",
            "--name",
            "morchy-e2e-postgres",
            "-p",
            f"{DB_PORT}:5432",
            "-e",
            f"POSTGRES_USER={DB_USER}",
            "-e",
            f"POSTGRES_PASSWORD={DB_PASSWORD}",
            "-e",
            f"POSTGRES_DB={DB_NAME}",
            "postgres:16-alpine",
        ],
        env=os.environ.copy() if "DOCKER_HOST" in os.environ else None,
    )

    time.sleep(5)
    wait_for_port(DB_HOST, DB_PORT, timeout=60)
    log("PostgreSQL ready")


def run_migrations():
    log("Running migrations...")
    conn_string = f"postgres://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}"
    env = os.environ.copy()
    env["GOOSE_DRIVER"] = "postgres"
    env["GOOSE_DBSTRING"] = conn_string
    run_cmd(["goose", "up", "-dir", str(MIGRATIONS_DIR)], env=env)
    log("Migrations done")


def clean_db():
    conn = db_conn()
    conn.autocommit = True
    cur = conn.cursor()
    for table in ["event", "lease", "spec", "workload", "goose_db_version"]:
        cur.execute(f"DROP TABLE IF EXISTS {table} CASCADE")
    cur.close()
    conn.close()


def start_controlplane():
    log("Starting control plane...")
    db_url = (
        f"postgresql://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}?sslmode=disable"
    )
    env = os.environ.copy()
    cp = subprocess.Popen(
        [str(BIN_DIR / "controlplane"), "--db", db_url, "--port", str(CP_PORT)],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        env=env,
    )
    wait_for_port("localhost", CP_PORT)
    log("Control plane ready")
    return cp


def start_agent(node_id, reserved_cpu=1000, reserved_ram=4096):
    log(f"Starting agent {node_id} with CPU={reserved_cpu}, RAM={reserved_ram}...")
    env = os.environ.copy()
    log_file = open(f"/tmp/agent_{node_id}.log", "w")
    agent = subprocess.Popen(
        [
            str(BIN_DIR / "agent"),
            "--controlplane", f"http://localhost:{CP_PORT}",
            "--node-id", node_id,
            "--reserved-cpu", str(reserved_cpu),
            "--reserved-ram", str(reserved_ram),
        ],
        stdout=log_file,
        stderr=subprocess.STDOUT,
        text=True,
        env=env,
    )
    time.sleep(2)

    if agent.poll() is not None:
        raise RuntimeError(f"Agent {node_id} failed to start")

    return agent


def create_workload(spec):
    resp = requests.post(f"{CP_API_URL}/workloads", json=spec, timeout=10)
    resp.raise_for_status()
    return resp.json()


def get_workload(wl_id):
    resp = requests.get(f"{CP_API_URL}/workloads/{wl_id}", timeout=10)
    resp.raise_for_status()
    return resp.json()


def list_workloads(status=None):
    params = {"status": status} if status else None
    resp = requests.get(f"{CP_API_URL}/workloads", params=params, timeout=10)
    resp.raise_for_status()
    return resp.json()


def get_lease(wl_id):
    conn = db_conn()
    cur = conn.cursor()
    cur.execute("SELECT node_id FROM lease WHERE workload_id = %s", (wl_id,))
    row = cur.fetchone()
    cur.close()
    conn.close()
    return {"node_id": row[0]} if row else None


def wait_for_status(wl_id, status, timeout=TIMEOUT):
    start = time.time()
    last_status = None
    check_interval = 5
    last_check = 0
    while time.time() - start < timeout:
        wl = get_workload(wl_id)
        last_status = wl["status"]
        if wl["status"] == status:
            return wl

        if time.time() - last_check >= check_interval:
            result = run_cmd(["docker", "ps", "-a", "--filter", f"name={wl_id}"], check=False, env=os.environ.copy())
            if result.stdout.strip():
                lines = result.stdout.strip().split('\n')
                if len(lines) > 1:
                    log(f"Container: {lines[1]}")
            last_check = time.time()

        time.sleep(POLL_INTERVAL)
    raise RuntimeError(f"Timeout waiting for status {status}, last status: {last_status}")


def wait_for_lease(wl_id, timeout=TIMEOUT):
    start = time.time()
    while time.time() - start < timeout:
        lease = get_lease(wl_id)
        if lease:
            return lease
        time.sleep(POLL_INTERVAL)
    raise RuntimeError("Timeout waiting for lease")


def wait_for_container_running(wl_id, timeout=TIMEOUT):
    start = time.time()
    while time.time() - start < timeout:
        result = run_cmd(["docker", "ps", "--filter", f"name={wl_id}", "--format", "{{.Status}}"], check=False, env=os.environ.copy())
        if result.stdout.strip() and "Up" in result.stdout:
            return True
        time.sleep(POLL_INTERVAL)
    raise RuntimeError(f"Timeout waiting for container to be running")

def test_basic_acquisition():
    log("\n=== Test: Basic Workload Acquisition ===")
    clean_db()
    run_migrations()

    cp = start_controlplane()
    agent = start_agent("agent-1")

    try:
        spec = {
            "image": "busybox:latest",
            "command": ["sh", "-c", "while true; do echo test; sleep 10; done"],
            "cpu": 50,
            "ram": 128,
            "env": {},
        }
        wl = create_workload(spec)
        assert wl["status"] == "new", f"Expected 'new', got {wl['status']}"
        log(f"Created workload {wl['id']}")

        wait_for_container_running(wl["id"])
        log("Container is running")

        wl = get_workload(wl["id"])
        assert wl["status"] in ["pending", "active"], f"Expected 'pending' or 'active', got {wl['status']}"
        log(f"Workload status: {wl['status']}")

        lease = wait_for_lease(wl["id"])
        assert lease["node_id"] == "agent-1", f"Expected agent-1, got {lease['node_id']}"
        log(f"Lease held by {lease['node_id']}")

        conn = db_conn()
        cur = conn.cursor()
        cur.execute("SELECT COUNT(*) FROM lease WHERE workload_id = %s", (wl["id"],))
        count = cur.fetchone()[0]
        cur.close()
        conn.close()
        assert count == 1, f"Expected 1 lease, got {count}"
        log("Test PASSED")
    except Exception as e:
        log(f"Checking agent logs...")
        try:
            with open("/tmp/agent_agent-1.log", "r") as f:
                lines = f.readlines()
                errors = [line for line in lines if "error" in line.lower() or "fail" in line.lower() or "panic" in line.lower()]
                if errors:
                    log(f"Agent errors: {''.join(errors[-10:])}")
                else:
                    log(f"Last 20 lines of agent log: {''.join(lines[-20:])}")
        except Exception:
            pass
        raise e
    finally:
        agent.terminate()
        agent.wait()
        cp.terminate()
        cp.wait()


def test_multiple_workloads():
    log("\n=== Test: Multiple Workloads ===")
    clean_db()
    run_migrations()

    cp = start_controlplane()
    agent = start_agent("agent-1")

    try:
        count = 3
        wls = []
        for i in range(count):
            spec = {
                "image": "busybox:latest",
                "command": ["sh", "-c", f"while true; do echo multi-{i}; sleep 10; done"],
                "cpu": 50,
                "ram": 128,
                "env": {},
            }
            wl = create_workload(spec)
            wls.append(wl)
            log(f"Created workload {wl['id']}")

        for wl in wls:
            wait_for_container_running(wl["id"])
            lease = wait_for_lease(wl["id"])
            assert lease["node_id"] == "agent-1"
            log(f"Workload {wl['id']} is running")

        log("Test PASSED")
    finally:
        agent.terminate()
        agent.wait()
        cp.terminate()
        cp.wait()


def test_exclusive_leases():
    log("\n=== Test: Exclusive Leases ===")
    clean_db()
    run_migrations()

    cp = start_controlplane()
    agents = [start_agent(f"agent-{i}") for i in range(3)]

    try:
        spec = {
            "image": "busybox:latest",
            "command": ["sh", "-c", "while true; do echo exclusive; sleep 10; done"],
            "cpu": 50,
            "ram": 128,
            "env": {},
        }
        wl = create_workload(spec)
        log(f"Created workload {wl['id']}")

        wait_for_container_running(wl["id"])
        lease = wait_for_lease(wl["id"])
        node_id = lease["node_id"]
        assert node_id in ["agent-0", "agent-1", "agent-2"], f"Unexpected node: {node_id}"
        log(f"Workload acquired by {node_id}")

        conn = db_conn()
        cur = conn.cursor()
        cur.execute("SELECT COUNT(*) FROM lease WHERE workload_id = %s", (wl["id"],))
        count = cur.fetchone()[0]
        cur.close()
        conn.close()
        assert count == 1, f"Expected 1 lease, got {count}"
        log("Test PASSED")
    finally:
        for a in agents:
            a.terminate()
            a.wait()
        cp.terminate()
        cp.wait()


def test_rescheduling():
    log("\n=== Test: Rescheduling on Agent Failure ===")
    clean_db()
    run_migrations()

    cp = start_controlplane()
    agent_a = start_agent("agent-a")
    agent_b = start_agent("agent-b")

    try:
        spec = {
            "image": "busybox:latest",
            "command": ["sh", "-c", "while true; do echo reschedule; sleep 10; done"],
            "cpu": 50,
            "ram": 128,
            "env": {},
        }
        wl = create_workload(spec)
        log(f"Created workload {wl['id']}")

        wait_for_container_running(wl["id"])
        lease = wait_for_lease(wl["id"])
        assert lease["node_id"] == "agent-a", f"Expected agent-a, got {lease['node_id']}"
        log("Workload initially on agent-a")

        log("Stopping agent-a...")
        agent_a.terminate()
        agent_a.wait()

        log("Waiting for rescheduling...")
        start = time.time()
        while time.time() - start < TIMEOUT:
            lease = get_lease(wl["id"])
            if lease and lease["node_id"] == "agent-b":
                break
            time.sleep(POLL_INTERVAL)
        else:
            raise RuntimeError("Rescheduling timeout")

        wait_for_container_running(wl["id"])
        wl = get_workload(wl["id"])
        assert wl["status"] in ["pending", "active"], f"Expected pending or active, got {wl['status']}"
        log("Workload rescheduled to agent-b")
        log("Test PASSED")
    finally:
        for a in [agent_a, agent_b]:
            if a.poll() is None:
                a.terminate()
                a.wait()
        cp.terminate()
        cp.wait()


def main():
    log("=" * 60)
    log("Starting E2E Tests")
    log("=" * 60)

    try:
        start_postgres()
        run_migrations()

        test_basic_acquisition()
        test_multiple_workloads()
        test_exclusive_leases()
        test_rescheduling()

        log("\n" + "=" * 60)
        log("ALL TESTS PASSED!")
        log("=" * 60)

    except Exception as e:
        log(f"\nTest failed: {e}")
        traceback.print_exc()
        sys.exit(1)
    finally:
        log("\nCleaning up...")
        result = run_cmd(["docker", "ps", "-q", "-f", "name=morchy-e2e-postgres"], check=False)
        if result.stdout.strip():
            run_cmd(["docker", "rm", "-f", result.stdout.strip().strip()], check=False)


if __name__ == "__main__":
    main()
