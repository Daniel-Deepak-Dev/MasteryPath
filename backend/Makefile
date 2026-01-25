.PHONY: swagger run build generate

# Generate Swagger documentation
swagger:
	go run github.com/swaggo/swag/cmd/swag init -g cmd/api/main.go

# Run the application
run:
	go run cmd/api/main.go

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run all go generate commands
generate:
	go generate ./...

# Seed the database with mock data
seed:
	go run scripts/seed.go
