swagger:
	swag init -g ./cmd/controlplane/main.go --parseInternal
	swagger-codegen generate -i docs/swagger.yaml -l go -o ./pkg/mctl/generated/controlplane.go

start-controlplane-dev: swagger
	go run cmd/controlplane/main.go --db "postgres://user:pass@localhost:5432/database?sslmode=disable"

start-agent-dev:
	go run cmd/agent/main.go --controlplane "http://localhost:8080" --reserved-ram 1024 --reserved-cpu 500 --node-id cf9c57a7-de96-4b3f-9a15-c44d4a57f5e1

build:
	go build -o bin/controlplane cmd/controlplane/main.go
	go build -o bin/agent cmd/agent/main.go
	go build -o bin/mctl cmd/mctl/main.go