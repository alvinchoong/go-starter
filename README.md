# Go Starter

Go Starter is a ready-to-use repo with modern tooling and best practices. It provides a solid foundation for building scalable Go applications following idiomatic Go practices.

## Dependencies and Tools

- **HTTP Router**:
  - [go-chi](https://github.com/go-chi/chi) - Lightweight, idiomatic and composable router for building Go HTTP services
  - [go-httphandler](https://github.com/alvinchoong/go-httphandler) - Idomatic HTTP request and response handling
- **Database**:
  - [jackc/pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit with advanced features and performance
  - [sqlc](https://github.com/sqlc-dev/sqlc) - Type-safe SQL code generation
  - [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations
- **Testing**: [stretchr/testify](https://github.com/stretchr/testify) - Toolkit with common assertions and mocks
- **Development**:
  - [golangci-lint](https://github.com/golangci/golangci-lint) - Fast Go linters runner
  - [air](https://github.com/cosmtrek/air) - Live reload for Go apps

## Project Structure

```plaintext
.
├── cmd/              # Application entry points
│   └── server/       # HTTP server implementation
├── internal/         # Private application code
│   ├── models/       # Generated database models
│   ├── mocks/        # Mock implementations for testing
│   └── pkg/          # Shared internal packages
├── database/         # Database related files
│   ├── migrations/   # SQL migrations
│   └── queries/      # SQL queries for sqlc
├── docs/             # Documentation
├── tools/            # Development tools
├── vendor/           # Vendored dependencies
└── build/            # Build artifacts
```

## Prerequisites

- [Go 1.23](https://go.dev/doc/go1.23) or later
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Getting Started

To use this template for your project:

1. Clone the repository

   ```bash
   git clone https://github.com/alvinchoong/go-starter.git your-project-name
   cd your-project-name
   ```

2. Replace all occurrences of `go-starter` with your module name

   ```bash
   # macOS/Linux
   # Replace in all Go files, go.mod, and Makefile
   find . -type f \( -name "*.go" -o -name "go.mod" -o -name "Makefile" \) -exec sed -i '' 's|go-starter|your-project-name|g' {} +
   ```

3. Initialize a new git repository

   ```bash
   # Remove existing git history
   rm -rf .git

   # Initialize new repository
   git init

   # Create initial commit
   git add .
   git commit -m "Initial commit"

   # Set up your remote
   git remote add origin https://github.com/username/your-project-name.git
   ```

4. Update project information
   - Update the LICENSE file accordingly
   - Modify this README to reflect your project's details

## Development

### Available Make Commands

- `make up`: Start PostgreSQL database with Docker Compose
- `make down`: Stop and remove Docker Compose services
- `make migrate`: Run database migrations
- `make sqlc`: Generate type-safe SQL code
- `make test`: Run tests with race detection
- `make server-run`: Run server with live reload
- `make server-build`: Build server binary
- `make server-docker-build`: Build server Docker image
- `make git-prepush-install`: Install git pre-push hook to run checks

## References

- [Effective Go](https://golang.org/doc/effective_go)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

## License

This project is licensed under the [MIT License](LICENSE).
