# Variables
APP_NAME=insightflow-users
export DB_URL?=postgres://postgres:postgres@localhost:5432/insightflow_users?sslmode=disable

# Build
.PHONY: build
build:
	go build -o bin/$(APP_NAME) ./cmd/server

# Run
.PHONY: run
run: build
	./bin/$(APP_NAME)

# Test
.PHONY: test
test:
	go test -v ./...

# Migrations
.PHONY: migrate-up
migrate-up:
	go run ./cmd/migrate/main.go up

.PHONY: migrate-down
migrate-down:
	go run ./cmd/migrate/main.go down

# Seeding
.PHONY: seed
seed:
	go run ./cmd/seed/main.go

.PHONY: reseed
reseed: migrate-down migrate-up seed

# Clean
.PHONY: clean
clean:
	rm -rf bin

# Initialize and run server
.PHONY: init
init: clean migrate-up seed run

# SQLC code generation
.PHONY: sqlc
sqlc:
	sqlc generate
