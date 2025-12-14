# Variables
APP_NAME=insightflow-users
DB_URL=sqlite3://users.db

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
reseed:
	rm -f users.db
	goose up
	go run ./cmd/seed/main.go

# Clean
.PHONY: clean
clean:
	rm -rf bin
	rm -f users.db

# Initialize and run server
.PHONY: init
init: clean migrate-up seed run