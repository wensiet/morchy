start-dev:
	swag init -g ./cmd/app/main.go --parseInternal
	go run cmd/app/main.go 