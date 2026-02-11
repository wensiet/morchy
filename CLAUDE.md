# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Morchy is a lightweight, stateless workload orchestration platform built around a pull-model based node agents. It orchestrates short-lived workloads with a minimal control plane and pull-based agents.

## Architecture

### Core Components

The system consists of four main components:

1. **Control Plane** (`cmd/controlplane/`, `pkg/controlplane/`)
   - Stateless API server that manages workload definitions
   - Stores workload state in PostgreSQL
   - Exposes REST API for workload CRUD operations, leasing, events, and edge configuration
   - Runs background tasks (e.g., lease expiration)
   - Uses Uber FX for dependency injection

2. **Agent** (`cmd/agent/`, `pkg/agent/`)
   - Pull-based worker deployed on runtime nodes
   - Fetches available workloads from control plane
   - Acquires leases on workloads to indicate ownership
   - Manages container lifecycle (create, monitor, reconcile, terminate)
   - Reports events back to control plane
   - Uses a supervisor pattern for managing periodic reconciliation tasks

3. **Edge** (`cmd/edge/`, `pkg/edge/`)
   - Network proxy component
   - Routes external traffic to running workloads
   - Fetches edge configuration from control plane

4. **mctl** (`cmd/mctl/`, `pkg/mctl/`)
   - CLI tool for interacting with the control plane
   - Supports creating, listing, and managing workloads
   - Accepts workload specs in YAML or JSON format

### Pull-Based Model

Unlike traditional push-based orchestrators:
- Agents **poll** the control plane for available workloads
- Agents **acquire leases** to claim workloads
- The control plane remains stateless about which agents are running
- This design improves resilience and simplifies scaling

### Layered Architecture

Each component follows a clean architecture pattern:
- **`domain/`**: Core business entities, value objects, and repository interfaces
- **`usecase/`**: Application-specific business rules and orchestration
- **`infrastructure/`**: External concerns (database, HTTP, configuration)
- **`implementation/`**: Concrete implementations of domain interfaces

## Development Commands

### Running Locally

```bash
# Start PostgreSQL database
docker-compose up -d

# Start control plane (requires database)
make start-controlplane-dev

# Start agent (in separate terminal)
make start-agent-dev

# Start edge proxy (in separate terminal)
make start-edge-dev
```

### Building

```bash
# Build all binaries to bin/
make build
```

### Code Generation

```bash
# Generate Swagger documentation and Go client
make swagger
```

This regenerates:
- `docs/swagger.yaml` - OpenAPI specification
- `pkg/mctl/generated/controlplane.go/` - Generated Go client

### Testing

```bash
# Run all tests
go test ./...

# Run tests in a specific package
go test ./pkg/controlplane/...
```

### Database Migrations

Migrations are located in `migrations/` and use the goose format:
- `20251029195659_init.sql` - Initial schema (workloads, specs, leases)
- `20251206095727_add_container.sql` - Container runtime support
- `20251211185006_add_event.sql` - Events table
- `20251214082334_extend_spec.sql` - Extended spec fields
- `20260105085927_add_net_config.sql` - Network configuration

Key tables:
- `workload` - Workload instances with status
- `spec` - Workload specifications (CPU, RAM, image, command, env, ports)
- `lease` - Agent leases on workloads (tracks which node is running what)
- `event` - System events from agents

## Key Concepts

### Workload Lifecycle

Workload statuses (defined in `pkg/controlplane/domain/workload/workload.go`):
- `new` - Created, waiting to be picked up
- `pending` - Agent is starting the workload
- `active` - Running successfully
- `stuck` - Failed to start after retries
- `failed` - Terminated with error
- `degraded` - Running but unhealthy

### Lease Management

- Agents acquire leases to claim workloads
- Leases have `updated_at` timestamps that agents refresh periodically
- Background task in control plane expires stale leases (updated > 30s ago)
- When a lease expires, the workload becomes available for other agents

### Workload Specifications

Workloads specify:
- **Image**: Docker image to run
- **CPU/RAM**: Resource limits in millicores and megabytes
- **Command**: Override container entrypoint
- **Env**: Environment variables
- **ContainerPort/HostPort**: Optional port mapping for networking

### Supervisor Pattern

The agent uses a generic supervisor pattern (`pkg/agent/infrastructure/supervisor.go`) to manage periodic tasks:
- Reconciliation loops that check workload health
- Lease renewal
- Each task runs at configurable intervals
- Tasks can be started/stopped dynamically

## Dependency Injection

All components use Uber FX for dependency injection. See `cmd/*/app/di.go` for wiring:
- Constructors create dependencies (logger, database, repositories, handlers)
- FX manages lifecycle (start/stop hooks)
- Background tasks are registered with lifecycle hooks

## API Structure

The control plane exposes REST endpoints:
- `POST /workloads` - Create workload
- `GET /workloads` - List workloads (with filters)
- `GET /workloads/:id` - Get workload by ID
- `DELETE /workloads/:id` - Delete workload
- `POST /workloads/:id/status` - Update workload status
- `POST /leases` - Acquire/release lease
- `POST /events` - Create event
- `GET /edges` - Get edge configurations

Full API spec: `docs/swagger.yaml`

## Configuration

### Control Plane
- `--db` flag: PostgreSQL connection string
- `--port` flag: HTTP server port (default varies)

### Agent
- `--controlplane` flag: Control plane URL
- `--node-id` flag: Unique node identifier
- `--reserved-ram` flag: RAM to reserve for system (MB)
- `--reserved-cpu` flag: CPU to reserve for system (millicores)

### Edge
- `--controlplane` flag: Control plane URL

### mctl
- `CONTROL_PLANE_URL` env var: Control plane URL (default: http://localhost:8080)

## Environment Variables

Check `internal/infrastructure` and `internal/domain` files for specific environment variable keys. Database connection is typically via `DATABASE_URL` or command-line flag.

## Workload Example

```yaml
# sandbox/workload.yaml
image: "busybox:latest"
command:
  - "sh"
  - "-c"
  - "while true; do echo \"hello\"; sleep 10; done"
cpu: 50
ram: 128
env: {}
```

Create with mctl:
```bash
./bin/mctl apply workload -f sandbox/workload.yaml
```

## Code Style

- DO NOT add single-line comments
- DO NOT add comments for anything

## Important Notes

- The control plane is **stateless** - it doesn't track agent health or maintain runtime state beyond what's in the database
- Agents are **responsible** for renewing leases and reporting workload status
- The system is designed for **short-lived workloads** - not long-running services
- Network routing through the edge component is optional and only needed for workloads with exposed ports
