# Morchy Codebase Guide for AI Agents

## Architecture Overview

**Morchy** is a **pull-based workload orchestration platform** with three core components:

1. **Control Plane** (`pkg/controlplane/`) - Stateless API server that manages workload definitions and leases
2. **Agents** (`pkg/agent/`) - Pull-based workers that poll the control plane for workloads and execute them via Docker
3. **Edge** (`pkg/edge/`) - Proxy component for agent communication

### Key Data Flow
- Agents call control plane API to pull workloads (pull-model, not push)
- Control plane stores workload specs in PostgreSQL (migrations in `migrations/`)
- Agents execute workloads in Docker and report status back
- Each component is independently deployable (stateless design)

## Project Structure & Patterns

### Layered Architecture Pattern
Each service follows **domain-driven design**:
```
pkg/{service}/
├── domain/          # Core entities & interfaces (no dependencies)
├── implementation/  # Concrete implementations
├── infrastructure/  # External integrations (DB, HTTP, Docker)
├── usecase/        # Business logic orchestrating domain
```

**Example**: [pkg/agent/](pkg/agent/) has:
- `domain/workload/` - Workload entity interfaces
- `implementation/repository/workload/` - In-memory workload storage
- `implementation/controlplane/` - HTTP client to control plane
- `usecase/handler.go` - Orchestrates the join/pull logic

### Dependency Injection with Uber Fx
All services use `go.uber.org/fx` for DI:
- `cmd/{service}/app/di.go` - Provides all dependencies as fx Providers
- Constructor functions create components: `func new{Component}(...) {Component}`
- [cmd/agent/app/di.go](cmd/agent/app/di.go) - Agent DI setup
- [cmd/controlplane/app/di.go](cmd/controlplane/app/di.go) - Control plane DI setup

**Pattern**: Add new components by defining provider functions in DI, fx automatically wires them.

### HTTP Framework (Gin)
- Control plane uses Gin router at `pkg/controlplane/infrastructure/router.go`
- Route handlers in `pkg/controlplane/implementation/gin.router/`
- Swagger/OpenAPI docs auto-generated from code comments at `docs/`

## Critical Workflows

### Local Development
```bash
make start-controlplane-dev    # Requires: DATABASE_URL env var (PostgreSQL)
make start-agent-dev           # Connects to http://localhost:8080 by default
make swagger                   # Regenerate OpenAPI docs
```

### Database Migrations
- Located in `migrations/` directory
- PostgreSQL is the primary datastore
- Run migrations before starting control plane (typically in deployment)

### Build & Deployment
```bash
make build                     # Produces bin/{controlplane,agent,mctl}
docker-compose up --build      # Full stack with PostgreSQL
```

## Key Conventions

### Workload Lifecycle States
- **Control Plane** (more states): `new` → `pending` → `active` → `failed|degraded|stuck|terminated`
- **Agent** (simpler): `new` → `active` → `terminated`
  
Implementation: [pkg/controlplane/domain/workload/workload.go](pkg/controlplane/domain/workload/workload.go) vs [pkg/agent/domain/workload/workload.go](pkg/agent/domain/workload/workload.go)

### Error Handling
Uses `samber/oops` for structured error wrapping. See `pkg/{service}/domain/errors.go` for domain-specific errors.

### Background Tasks
Control plane runs periodic tasks (e.g., lease expiration) via `BackgroundTaskRunner` in `pkg/controlplane/infrastructure/background.go`. Register tasks in `di.go`.

### Agent Pull Loop
Agent runs continuous join logic every 10 seconds in `runLoop()` function in [cmd/agent/app/di.go](cmd/agent/app/di.go):
1. `LoadCurrentState()` - Sync Docker containers with local state
2. `ApplyWorkloadJoin()` - Pull new workloads from control plane

## Important Files for Context

- [go.mod](go.mod) - Key dependencies: Gin, pgx, Docker client, fx, Zap logging
- [cmd/](cmd/) - Service entrypoints; each has `app/{command,config,di}.go`
- [pkg/controlplane/implementation/gin.router/](pkg/controlplane/implementation/gin.router/) - All HTTP endpoints
- [pkg/agent/implementation/controlplane/client.go](pkg/agent/implementation/controlplane/client.go) - Agent→ControlPlane API calls
- [pkg/runtime/client.go](pkg/runtime/client.go) - Docker container execution wrapper

## Testing & Validation

- Run tests: `go test ./...`
- Lint: Standard Go linting recommended (none configured yet)
- Manual testing: Use `sandbox/` directory for test configs (YAML workloads, JSON container definitions)

## Common Tasks for Agents

- **Adding new workload fields**: Update [pkg/controlplane/domain/workload/workload.go](pkg/controlplane/domain/workload/workload.go) struct, add DB migration, update repository queries
- **Adding new API endpoint**: Create handler in `pkg/controlplane/implementation/gin.router/`, add route in `router.go`, add Swagger comments
- **Adding background task**: Define in control plane DI (`RegisterTask()`), implement handler in usecase
- **Adding agent feature**: Update `pkg/agent/usecase/handler.go` to add join logic, then update DI if new components needed
