# Users Service - Workshop 3 - System's Architecture

This repository contains the users service used by the **InsightFlow** system from the third workshop of the subject "Arquitectura de Sistemas" at Universidad Católica del Norte. Below are the tools used and how to setup this project locally.

## Pre-requisites

- [Go](https://golang.org/dl/) (version 1.21+)
- [PostgreSQL](https://www.postgresql.org/) (version 15+)
- [Git](https://git-scm.com/) (version 2.49.0)
- [Docker or Docker Desktop](https://docs.docker.com/)
- [SQLC](https://sqlc.dev/) (optional, for code generation)

**Note**: This project can be setup either using the first three pre-requisites or using only Docker. It is recommended to use just **Docker**.

## Installation and configuration

1. **Clone the repository**

```bash
git clone https://github.com/Taller-3-Arq-de-Sistemas/insightflow-users.git
```

2. **Navigate to the project directory**

```bash
cd insightflow-users
```

## Setup using Docker

3.1. **Create a `.env` file using the example environment variables file and fill its values**

```bash
cp .env.example .env
```

In the `.env` file, you can replace:

- `PORT` with the port that the app uses to expose the server.
- `DB_URL` with the database connection string.
- `JWT_SECRET` with the JWT secret that you want to use.

Once you have replaced everything, save the changes and move on to the next step.

3.2. **Build and run the project using docker compose**

```bash
docker compose -f docker-compose.yaml -f docker-compose.override.yaml up --build -d
```

This will start both the PostgreSQL database and the application. The web app will be running on port **8080**.

## Setup without Docker

4.1. **Start a PostgreSQL instance**

You need a PostgreSQL database running. You can use Docker to run just the PostgreSQL instance:

```bash
docker compose -f docker-compose.yaml -f docker-compose.override.yaml up -d postgres
```

If you already have Postgres on port 5432 locally, adjust the compose port mapping to avoid conflicts.

4.2. **Create a `.env` file**

```bash
cp .env.example .env
```

Update the `DB_URL` to match your local PostgreSQL configuration:

```
DB_URL=postgres://postgres:postgres@localhost:5432/insightflow_users?sslmode=disable
```

4.3. **Run database migrations**

```bash
make migrate-up
```

4.4. **Seed the database**

```bash
make seed
```

4.5. **Run the project**

```bash
make run
```

The server will start on the configured port (default: **8080**). Access the API via http://localhost:8080.

**Alternative**: You can run all setup steps at once using:

```bash
make init
```

## Available Make commands

| Command | Description |
|---------|-------------|
| `make build` | Build the application binary |
| `make run` | Build and run the application |
| `make test` | Run all tests |
| `make migrate-up` | Run database migrations |
| `make migrate-down` | Rollback database migrations |
| `make seed` | Seed the database with initial data |
| `make reseed` | Reset and reseed the database |
| `make clean` | Remove build artifacts |
| `make init` | Clean, migrate, seed, and run the server |
| `make sqlc` | Regenerate SQLC code |

## Data seeder

The seeder is run automatically with `make seed` or `make init`.

The seeder contains:

- 3 roles (admin, user, editor)
- 251 users (1 administrator and 250 randomly generated users)

**Default admin credentials:**
- Email: `jairo.calcina@admin.com`
- Password: `AReallyGoodP4ssw0rd`

## Operations

The operations that this application exposes are 9 and are separated in two modules: the users module and the authentication module. Below is a detailed overview of these operations:

### Authentication module

Base path: `/api/v1/auth`

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/login` | Login with email and password | No |
| POST | `/register` | Register a new user | No |
| POST | `/logout` | Logout and invalidate token | Yes |
| POST | `/validate` | Validate a JWT token | No |

#### Login

- **URI**: `/api/v1/auth/login`
- **Method**: POST
- **Body**:
  - `email`: Email of the user (required)
  - `password`: Password of the user (required)
- **Response**: `{ "token": "jwt_token_here" }`

#### Register

- **URI**: `/api/v1/auth/register`
- **Method**: POST
- **Body**:
  - `name`: Name of the user (required)
  - `last_names`: Last names of the user (required)
  - `email`: Email of the user (required)
  - `username`: Username (required)
  - `password`: Password with at least 8 characters (required)
  - `birth_date`: Birth date in YYYY-MM-DD format (required)
  - `address`: Address of the user (required)
  - `phone`: Phone number, 9-15 digits (required)
- **Response**: `{ "token": "jwt_token_here" }`

#### Validate Token

- **URI**: `/api/v1/auth/validate`
- **Method**: POST
- **Body**:
  - `token`: JWT token to validate (required)
- **Response**: `{ "id": "user_id", "role": "user_role", "exp": expiration_timestamp }`

### Users module

Base path: `/api/v1/users`

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/` | Create a new user | Yes (Admin) |
| GET | `/` | List all active users | Yes |
| GET | `/{id}` | Get user by ID | No |
| PATCH | `/{id}` | Update user (own profile only) | Yes |
| DELETE | `/{id}` | Delete user (soft delete) | Yes (Admin) |

#### Create User (Admin only)

- **URI**: `/api/v1/users`
- **Method**: POST
- **Body**:
  - `name`: Name of the user (required)
  - `last_names`: Last names of the user (required)
  - `email`: Email of the user (required)
  - `username`: Username (required)
  - `password`: Password with at least 8 characters (required)
  - `birth_date`: Birth date in YYYY-MM-DD format (required)
  - `address`: Address of the user (required)
  - `phone`: Phone number (required)
  - `status`: User status (`active` or `inactive`)
  - `role`: User role (`admin`, `user`, or `editor`)
- **Response**: User object

#### List Users

- **URI**: `/api/v1/users`
- **Method**: GET
- **Response**: Array of active users

#### Get User by ID

- **URI**: `/api/v1/users/{id}`
- **Method**: GET
- **Response**: User object

#### Update User

- **URI**: `/api/v1/users/{id}`
- **Method**: PATCH
- **Body**:
  - `full_name`: Full name (first name + last names)
  - `username`: New username
- **Response**: Updated user object

#### Delete User (Admin only)

- **URI**: `/api/v1/users/{id}`
- **Method**: DELETE
- **Response**: 200 OK

## Health endpoints

The service exposes health endpoints:

- `/`: Returns "Hello World"
- `/health`: Returns "Everything is OK"

## Service architecture

This service is part of the InsightFlow SOA system. It is composed by a server made in Go and a PostgreSQL database.

The database schema is composed by 3 tables:
- `role`: Stores user roles (admin, user, editor)
- `users`: Stores user information
- `token_blacklist`: Stores invalidated JWT tokens

## Design patterns applied

The service uses several patterns to keep the codebase modular, testable, and maintainable:

### Layered Architecture (Handler-Service-Repository)

The application follows a three-layer architecture that separates concerns:

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Layer                           │
│  ┌─────────────────┐         ┌─────────────────┐       │
│  │  users/handler  │         │  auth/handler   │       │
│  └────────┬────────┘         └────────┬────────┘       │
├───────────┼───────────────────────────┼─────────────────┤
│           ▼         Business Layer    ▼                 │
│  ┌─────────────────┐         ┌─────────────────┐       │
│  │  users/service  │         │  auth/service   │       │
│  └────────┬────────┘         └────────┬────────┘       │
├───────────┼───────────────────────────┼─────────────────┤
│           ▼          Data Layer       ▼                 │
│  ┌──────────────────────────────────────────────┐      │
│  │          postgres/sqlc (Repository)          │      │
│  └──────────────────────────────────────────────┘      │
└─────────────────────────────────────────────────────────┘
```

- **Handlers**: Receive HTTP requests, validate input, call services, and return responses
- **Services**: Contain business logic, orchestrate operations, and enforce rules
- **Repository**: Handle database operations via SQLC-generated code

### Repository Pattern

Database access is abstracted through the `Querier` interface generated by SQLC. This provides:

- Type-safe database queries
- Compile-time SQL validation
- Easy mocking for unit tests

```go
// Generated interface that services depend on
type Querier interface {
    CreateUser(ctx context.Context, arg CreateUserParams) (string, error)
    FindUserById(ctx context.Context, id string) (FindUserByIdRow, error)
    // ... other methods
}
```

### Dependency Injection

All components receive their dependencies through constructors rather than creating them internally:

```go
// Services receive repository via constructor
func NewService(repo repository.Querier) *svc {
    return &svc{repo: repo}
}

// Handlers receive services via constructor
func NewHandler(service *svc) *Handler {
    return &Handler{service: service}
}
```

This enables:
- Easy testing with mock dependencies
- Loose coupling between components
- Centralized dependency wiring in `main.go`

### Middleware Pattern

Cross-cutting concerns are implemented as Chi middleware that wraps handlers:

- **AuthMiddleware**: Validates JWT tokens and extracts user context
- **AdminMiddleware**: Enforces admin-only access on protected routes
- **Built-in middleware**: Request logging, panic recovery, timeouts, CORS

```go
r.Route("/users", func(r chi.Router) {
    r.With(authMiddleware, adminMiddleware).Post("/", handler.CreateUser)
    r.With(authMiddleware).Get("/", handler.ListUsers)
})
```

### DTO Pattern (Data Transfer Objects)

Request and response types are separate from domain/database models:

- **Request DTOs**: `CreateUserParams`, `LoginParams`, `RegisterParams`
- **Response DTOs**: SQLC-generated row types like `FindUserByIdRow`, `ListUsersRow`

This decouples the API contract from internal data structures and allows independent evolution.

### Adapter Pattern

The `internal/adapters/postgres` package acts as an adapter between the application and PostgreSQL:

- Encapsulates database connection management
- Provides connection pooling configuration
- Isolates database-specific code from business logic

```go
// Database adapter with connection pooling
func New(dsn string) (*DB, error) {
    db, err := sql.Open("pgx", dsn)
    // Configure pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    return &DB{db}, nil
}
```

### Context Pattern

User authentication data flows through the request context:

```go
// Set user info in context after authentication
ctx = SetUserContext(ctx, userID, role)

// Retrieve user info in handlers
userID, ok := auth.GetUserID(r.Context())
role, ok := auth.GetRole(r.Context())
```

This avoids passing authentication data through function parameters and keeps handlers clean.

## Tech stack

- **Language**: Go 1.21+
- **Router**: Chi
- **Database**: PostgreSQL
- **ORM/Query Builder**: SQLC
- **Migrations**: Goose
- **Authentication**: JWT (golang-jwt)
- **Validation**: go-playground/validator

## Project structure

```
.
├── cmd/
│   ├── server/      # Main application entry point
│   ├── migrate/     # Database migration tool
│   └── seed/        # Database seeder
├── config/          # Configuration loading
├── internal/
│   ├── adapters/
│   │   └── postgres/  # Database adapter and queries
│   ├── api/           # HTTP response helpers
│   ├── auth/          # Authentication service
│   ├── json/          # JSON utilities
│   ├── users/         # Users service
│   └── validator/     # Input validation
├── seeders/         # Seed data files
├── Dockerfile
├── docker-compose.yaml
└── Makefile
```

## Author

- [@Jairo Calcina](https://github.com/Broukt)
