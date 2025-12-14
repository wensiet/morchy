swagger:
	swag init -g ./cmd/controlplane/main.go --parseInternal

start-controlplane-dev: swagger
	go run cmd/controlplane/main.go --db "postgres://user:pass@localhost:5432/database?sslmode=disable"

start-agent-dev:
	go run cmd/agent/main.go --controlplane "http://localhost:8080" --reserved-ram 1024 --reserved-cpu 500

build:
	go build -o bin/controlplane cmd/controlplane/main.go
	go build -o bin/agent cmd/agent/main.go