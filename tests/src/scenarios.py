import asyncio
import logging
import time
from abc import ABC, abstractmethod
from typing import Any

import httpx

from src.settings import AgentSettings, ControlPlaneSettings
from src.ssh_client import SSHClient

logger = logging.getLogger(__name__)


class BaseScenario(ABC):
    def __init__(
        self,
        cp_settings: ControlPlaneSettings,
        agent_settings: AgentSettings,
    ):
        self._httpx_client = httpx.AsyncClient(timeout=30.0)
        self._cp_settings = cp_settings

        self._agent_ssh_clients = {}
        for agent_host in agent_settings.hosts:
            self._agent_ssh_clients[agent_host] = SSHClient(
                host=agent_host,
                username="root",
                key_filename=agent_settings.ssh_key_file_path,
            )

        super().__init__()

    async def _api_request(
        self, method: str, endpoint: str, data: dict | None = None
    ) -> tuple[int, dict[str, Any]]:
        url = f"{self._cp_settings.url.rstrip('/')}/{endpoint.lstrip('/')}"
        headers = {"Content-Type": "application/json"}

        try:
            if data:
                response = await self._httpx_client.request(
                    method, url, json=data, headers=headers
                )
            else:
                response = await self._httpx_client.request(
                    method, url, headers=headers
                )

            return response.status_code, response.json()
        except httpx.HTTPStatusError as e:
            return (
                e.response.status_code,
                e.response.json() if e.response.content else {},
            )
        except Exception as e:
            return 0, {"error": str(e)}

    async def _cleanup_all_workloads(self) -> None:
        code, workloads = await self._api_request("GET", "/api/v1/workloads")

        if code == 200 and workloads:
            for workload in workloads:
                workload_id = workload.get("id")
                if workload_id:
                    await self._api_request(
                        "DELETE", f"/api/v1/workloads/{workload_id}"
                    )

    @abstractmethod
    async def run(self) -> None:
        raise NotImplementedError


class BasicWorkloadAcquisition(BaseScenario):
    async def run(self) -> None:
        """
        Ensure that there is no workloads. If some - delete them.
        Apply some basic workload.
        Wait for it to become running.
        Verify that it is really running with ssh and docker ps on agents.
        """
        workload_spec = {
            "image": "busybox:latest",
            "command": ["sh", "-c", "while true; do echo hello; sleep 10; done"],
            "cpu": 50,
            "ram": 128,
            "env": {},
        }

        await self._cleanup_all_workloads()

        await asyncio.sleep(1)

        code, response = await self._api_request(
            "POST", "/api/v1/workloads", workload_spec
        )

        if code != 201 or "id" not in response:
            raise RuntimeError(f"Failed to create workload: {response}")

        workload_id = response["id"]

        logger.info(f"Workload {workload_id} applied")

        start_time = time.time()
        timeout = 90

        while time.time() - start_time < timeout:
            code, wl = await self._api_request(
                "GET", f"/api/v1/workloads/{workload_id}"
            )

            if code == 200 and wl.get("status") == "active":
                break

            await asyncio.sleep(0.5)
        else:
            raise RuntimeError(
                f"Workload {workload_id} did not reach active status within {timeout}s"
            )

        logger.info("Workload status transitioned to active")

        workload_running = False

        for agent_host, ssh_client in self._agent_ssh_clients.items():
            result = ssh_client.run("docker ps --format '{{.Names}}'")

            if workload_id in result.stdout:
                workload_running = True
                logger.info(f"Verified that workload ran on agent {agent_host}")
                break

        if not workload_running:
            raise RuntimeError(f"Workload {workload_id} is not running on any agent")


class MultipleWorkloads(BaseScenario):
    async def run(self) -> None:
        """
        Ensure that there is no workloads. If some - delete them.
        Apply two workloads.
        Wait for both to become running.
        Verify that both are really running with ssh and docker ps on agents.
        """
        workload_specs = [
            {
                "image": "busybox:latest",
                "command": [
                    "sh",
                    "-c",
                    "while true; do echo workload-1; sleep 10; done",
                ],
                "cpu": 50,
                "ram": 128,
                "env": {},
            },
            {
                "image": "busybox:latest",
                "command": [
                    "sh",
                    "-c",
                    "while true; do echo workload-2; sleep 10; done",
                ],
                "cpu": 50,
                "ram": 128,
                "env": {},
            },
        ]

        await self._cleanup_all_workloads()
        await asyncio.sleep(1)

        workload_ids = []

        for spec in workload_specs:
            code, response = await self._api_request("POST", "/api/v1/workloads", spec)

            if code != 201 or "id" not in response:
                raise RuntimeError(f"Failed to create workload: {response}")

            workload_id = response["id"]
            workload_ids.append(workload_id)
            logger.info(f"Workload {workload_id} applied")

        start_time = time.time()
        timeout = 90

        while time.time() - start_time < timeout:
            all_active = True

            for workload_id in workload_ids:
                code, wl = await self._api_request(
                    "GET", f"/api/v1/workloads/{workload_id}"
                )

                if code != 200 or wl.get("status") != "active":
                    all_active = False
                    break

            if all_active:
                break

            await asyncio.sleep(0.5)
        else:
            raise RuntimeError(
                f"Workloads did not reach active status within {timeout}s"
            )

        logger.info("All workloads status transitioned to active")

        for workload_id in workload_ids:
            workload_running = False

            for agent_host, ssh_client in self._agent_ssh_clients.items():
                result = ssh_client.run("docker ps --format '{{.Names}}'")

                if workload_id in result.stdout:
                    workload_running = True
                    logger.info(
                        f"Verified that workload {workload_id} ran on agent {agent_host}"
                    )
                    break

            if not workload_running:
                raise RuntimeError(
                    f"Workload {workload_id} is not running on any agent"
                )

        logger.info("Verified that all workloads ran on docker socket")


class Rescheduling(BaseScenario):
    async def run(self) -> None:
        """
        Ensure that there is no workloads. If some - delete them.
        Apply a workload.
        Wait for it to become active and identify which agent acquired it.
        Disable (stop) the agent running the workload.
        Wait for the workload to be acquired by the other agent and become active again.
        """
        if len(self._agent_ssh_clients) < 2:
            raise RuntimeError("Rescheduling scenario required at least 2 agents")

        workload_spec = {
            "image": "busybox:latest",
            "command": ["sh", "-c", "while true; do echo hello; sleep 10; done"],
            "cpu": 50,
            "ram": 128,
            "env": {},
        }

        await self._cleanup_all_workloads()
        await asyncio.sleep(1)

        code, response = await self._api_request(
            "POST", "/api/v1/workloads", workload_spec
        )

        if code != 201 or "id" not in response:
            raise RuntimeError(f"Failed to create workload: {response}")

        workload_id = response["id"]
        logger.info(f"Workload {workload_id} applied")

        start_time = time.time()
        timeout = 90

        while time.time() - start_time < timeout:
            code, wl = await self._api_request(
                "GET", f"/api/v1/workloads/{workload_id}"
            )

            if code == 200 and wl.get("status") == "active":
                break

            await asyncio.sleep(0.5)
        else:
            raise RuntimeError(
                f"Workload {workload_id} did not reach active status within {timeout}s"
            )

        logger.info("Workload status transitioned to active")

        primary_agent_host = None
        for agent_host, ssh_client in self._agent_ssh_clients.items():
            result = ssh_client.run("docker ps --format '{{.Names}}'")

            if workload_id in result.stdout:
                primary_agent_host = agent_host
                logger.info(f"Workload acquired by agent {agent_host}")
                break

        if not primary_agent_host:
            raise RuntimeError(f"Workload {workload_id} did not run on any agent")

        code, lease = await self._api_request(
            "GET", f"/api/v1/workloads/{workload_id}/lease"
        )

        if code != 200 or "node_id" not in lease:
            raise RuntimeError(f"Failed to get lease: {lease}")

        original_node_id = lease["node_id"]
        logger.info(f"Original lease acquired by node {original_node_id}")

        primary_ssh_client = self._agent_ssh_clients[primary_agent_host]
        primary_ssh_client.run("systemctl stop morchy-agent")
        logger.info(f"Stopped morchy-agent service on {primary_agent_host}")

        try:
            start_time = time.time()
            timeout = 60

            while time.time() - start_time < timeout:
                code, lease = await self._api_request(
                    "GET", f"/api/v1/workloads/{workload_id}/lease"
                )

                if code == 200:
                    current_node_id = lease.get("node_id")
                    if current_node_id != original_node_id:
                        logger.info(f"Lease acquired by new node {current_node_id}")
                        break
                elif code == 404:
                    logger.info("Lease released")
                    break

                await asyncio.sleep(0.5)
            else:
                raise RuntimeError(f"Lease did not change or expire within {timeout}s")

            logger.info("Waited for lease to be released or acquired by another node")

            start_time = time.time()
            timeout = 90

            while time.time() - start_time < timeout:
                code, wl = await self._api_request(
                    "GET", f"/api/v1/workloads/{workload_id}"
                )

                if code == 200 and wl.get("status") == "active":
                    break

                await asyncio.sleep(0.5)
            else:
                raise RuntimeError(
                    f"Workload {workload_id} did not reach active status after rescheduling within {timeout}s"
                )

            logger.info("Workload status transitioned to active after rescheduling")

        finally:
            await asyncio.sleep(2)
            primary_ssh_client.run("systemctl start morchy-agent")
            logger.info("Restarted stopped agent")

        workload_running = False
        running_on_agent = None

        for agent_host, ssh_client in self._agent_ssh_clients.items():
            if agent_host == primary_agent_host:
                continue

            result = ssh_client.run("docker ps --format '{{.Names}}'")

            if workload_id in result.stdout:
                workload_running = True
                running_on_agent = agent_host
                logger.info(f"Workload ran on agent {agent_host}")
                break

        if not workload_running:
            raise RuntimeError(
                f"Workload {workload_id} did not run on another agent after rescheduling"
            )

        await asyncio.sleep(2)

        result = primary_ssh_client.run("docker ps --format '{{.Names}}'")
        if workload_id in result.stdout:
            logger.info(
                f"Workload is still running on older agent {primary_agent_host}"
            )

        logger.info(
            f"Verified that workload rescheduled from {primary_agent_host} to {running_on_agent}"
        )


class ConcurrentWorkloads(BaseScenario):
    async def run(self) -> None:
        """
        Ensure that there is no workloads. If some - delete them.
        Disable all agents except one.
        Apply two workloads.
        Wait for them both to become active.
        """

        workload_specs = [
            {
                "image": "busybox:latest",
                "command": [
                    "sh",
                    "-c",
                    "while true; do echo workload-1; sleep 10; done",
                ],
                "cpu": 50,
                "ram": 128,
                "env": {},
            },
            {
                "image": "busybox:latest",
                "command": [
                    "sh",
                    "-c",
                    "while true; do echo workload-2; sleep 10; done",
                ],
                "cpu": 50,
                "ram": 128,
                "env": {},
            },
        ]

        await self._cleanup_all_workloads()
        await asyncio.sleep(1)

        target_agent_host = None
        for host, ssh_client in self._agent_ssh_clients.items():
            if target_agent_host is None:
                target_agent_host = host
                continue

            ssh_client.run("systemctl stop morchy-agent")

        logger.info(f"Disabled all agents except {target_agent_host}")

        try:
            workload_ids = []

            for spec in workload_specs:
                code, response = await self._api_request(
                    "POST", "/api/v1/workloads", spec
                )

                if code != 201 or "id" not in response:
                    raise RuntimeError(f"Failed to create workload: {response}")

                workload_id = response["id"]
                workload_ids.append(workload_id)
                logger.info(f"Workload {workload_id} applied")

            start_time = time.time()
            timeout = 90

            while time.time() - start_time < timeout:
                all_active = True

                for workload_id in workload_ids:
                    code, wl = await self._api_request(
                        "GET", f"/api/v1/workloads/{workload_id}"
                    )

                    if code != 200 or wl.get("status") != "active":
                        all_active = False
                        break

                if all_active:
                    break

                await asyncio.sleep(0.5)
            else:
                raise RuntimeError(
                    f"Workloads did not reach active status within {timeout}s"
                )

            logger.info("All workloads status transitioned to active")

            result = self._agent_ssh_clients[target_agent_host].run(
                "docker ps --format '{{.Names}}'"
            )

            for workload_id in workload_ids:
                if workload_id not in result.stdout:
                    raise RuntimeError(
                        f"Workload {workload_id} is not running on target agent"
                    )
        finally:
            for host, ssh_client in self._agent_ssh_clients.items():
                if host == target_agent_host:
                    continue

                ssh_client.run("systemctl start morchy-agent")
                logger.info(f"Started agent {host}")
