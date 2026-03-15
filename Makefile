swagger:
	swag init -g ./cmd/controlplane/main.go --parseInternal
	swagger-codegen generate -i docs/swagger.yaml -l go -o ./pkg/mctl/generated/controlplane.go

start-controlplane-dev: swagger
	go run cmd/controlplane/main.go --db "postgres://user:pass@localhost:5432/database?sslmode=disable"

start-agent-dev:
	go run cmd/agent/main.go --controlplane "http://localhost:8080" --reserved-ram 1024 --reserved-cpu 500 --node-id cf9c57a7-de96-4b3f-9a15-c44d4a57f5e1

start-edge-dev:
	go run cmd/edge/main.go --controlplane "http://localhost:8080"

build:
	go build -o bin/controlplane cmd/controlplane/main.go
	go build -o bin/agent cmd/agent/main.go
	go build -o bin/mctl cmd/mctl/main.go
	go build -o bin/edge cmd/edge/main.go

test:
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test
	go tool cover -html=coverage.out -o coverage.html

test-coverage-clean:
	go test -v -race -coverprofile=coverage.out ./pkg/...
	go tool cover -func=coverage.out | grep -v mocks | grep -v testutil
	go tool cover -html=coverage.out -o coverage.html

test-coverage-report:
	go tool cover -func=coverage.out | grep -v mocks | grep -v testutil

test-short:
	go test -short -v ./...

mocks:
	mockery --config .mockery.yaml