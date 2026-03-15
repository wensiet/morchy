# E2E Tests

Simple end-to-end tests for the Morchy orchestration system.

## Prerequisites

1. Build the binaries:
   ```bash
   make build
   ```

2. Install Python dependencies:
   ```bash
   pip install requests psycopg2-binary
   ```

3. Install goose for migrations:
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

4. Ensure Docker is running (for PostgreSQL container)

## Running Tests

```bash
cd tests/e2e
python run_e2e.py
```

## Test Scenarios

1. **Basic Workload Acquisition**: Single agent creates and acquires a workload
2. **Multiple Workloads**: Single agent handles multiple workloads simultaneously
3. **Exclusive Leases**: Multiple agents compete for single workload, only one gets lease
4. **Rescheduling**: Workload reschedules to another agent when primary agent fails

## Cleanup

The tests automatically clean up Docker containers on exit. If cleanup fails:
```bash
docker rm -f morchy-e2e-postgres
```
