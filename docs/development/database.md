# Database Operations Guide

This guide covers how to work with database migrations and queries in the project.

## Database Migrations

### Creating Migration Files

Migration files help you manage database schema changes in a version-controlled way. We use [golang-migrate](https://github.com/golang-migrate/migrate) for handling migrations.

You can create migration files manually following the naming convention, or use the golang-migrate CLI:

#### Using golang-migrate CLI

```bash
migrate create -ext sql -dir database/migrations -seq create_users_table
```

This will create two files:

- `{version}_create_users_table.up.sql`
- `{version}_create_users_table.down.sql`

#### Creating Migration Files Manually

1. Migration files should be created in the `database/migrations` directory
2. Files follow the naming convention: `{version}_{description}.{up|down}.sql`
   - Example: `000001_create_users_table.up.sql`
   - Example: `000001_create_users_table.down.sql`

### Running Migrations

To apply migrations:

```bash
make migrate  # applies all pending migrations
# or specify number of migrations
make migrate STEP=1  # applies 1 migration
```

To rollback migrations:

```bash
make migrate CMD=down  # rolls back all migrations
# or specify number of migrations
make migrate CMD=down STEP=2  # rolls back 2 migrations
```

To access the database console:

```bash
make db-console
```

## Working with Database Queries

We use [sqlc](https://sqlc.dev/) to generate type-safe Go code from SQL queries.

### Creating New Queries

1. Create or modify query files in the `database/queries` directory
2. Each query file should focus on a specific database entity or domain concept
3. Write your SQL queries with sqlc annotations

For detailed query syntax and annotations, refer to the [sqlc query documentation](https://docs.sqlc.dev/en/latest/reference/query-syntax.html).

Example query file (`database/queries/users.sql`):

```sql
-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (email)
VALUES ($1)
RETURNING *;
```

### Generating Query Code

After creating or modifying queries, generate the Go code:

```bash
make sqlc
```

This will create:

- Go structs for your database tables
- Type-safe methods for your queries
- Interface definitions for your queries, useful for mocking and unit testing

### Using Generated Queries

The generated code can be used in your services:

```go
import "your-project/service"

func (s *Service) GetUser(ctx context.Context, db models.DBTX, id uuid.UUID) (models.User, error) {
    return s.queries.GetUser(ctx, db, id)
}
```

## Best Practices

1. Always create both up and down migrations
2. Keep migrations atomic
3. Test migrations on a development database before applying to production
4. Document complex queries with comments
