start-controlplane-dev:
	swag init -g ./cmd/app/main.go --parseInternal
	go run cmd/controlplane/main.go

start-agent-dev:
	go run cmd/agent/main.go

build:
	go build -o bin/controlplane cmd/controlplane/main.go
	go build -o bin/agent cmd/agent/main.go