# Helper script to generate Swagger documentation
go run github.com/swaggo/swag/cmd/swag init -g cmd/api/main.go
Write-Host "Swagger documentation generated successfully!" -ForegroundColor Green
