SHELL = /bin/bash -u -e -o pipefail

# `make` applies env vars from `.env`
include .env
export

git-prepush-install:
	@echo "#!/bin/sh" > .git/hooks/pre-push
	@echo "make git-prepush" >> .git/hooks/pre-push
	@chmod +x .git/hooks/pre-push
	@echo "Git pre-push hook set up to run 'make git-prepush'"

git-prepush: lint

up:
	docker-compose up -d --remove-orphans

down:
	docker-compose down

# migrate target supports:
# make migrate (applies all pending migrations)
# make migrate CMD=up STEP=1 (applies 1 pending migrations)
# make migrate CMD=down STEP=2 (rolls back 2 migrations)
migrate: CMD=up
migrate: STEP=
migrate:
	which migrate || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@while ! pg_isready -q -d $(DATABASE_URL); do \
		echo "Waiting for PostgreSQL to be available..."; \
		sleep 1; \
	done
	migrate -path ./database/migrations -database "$(DATABASE_URL)?sslmode=disable" $(CMD) $(STEP)

db-console:
	psql $(DATABASE_URL)

sqlc:
	sqlc version | grep v1.28.0 || go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0
	sqlc generate

test:
	go test -v -race ./...

GIT_VERSION ?= $(shell git describe --tags --always)
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILDINFO_PKG=go-starter/internal/pkg/buildinfo
server-build:
	go build -o build/server -trimpath \
		-ldflags "-X $(BUILDINFO_PKG).Version=$(GIT_VERSION) \
		-X $(BUILDINFO_PKG).BuildTime=$(BUILD_TIME)" \
		cmd/server/main.go

server-docker-build:
	docker buildx build \
		--platform=linux/arm64 \
		-t go-starter:server \
		-f cmd/server/Dockerfile .

server-run:
	which air || go install github.com/air-verse/air@latest
	air --build.delay=1000 \
		--build.cmd "make server-build" \
		--build.bin "./build/server" \
		--build.include_ext "go" \
		--build.exclude_dir "vendor" \
		--build.exclude_regex ".*_test.go"

lint:
	go mod tidy
	if ! git diff --quiet go.mod go.sum; then \
		printf "$(RED)There are changes to the go.mod & go.sum files$(NORMAL)\n"; \
		exit 1; \
	fi
	gofumpt -version | grep v0.7.0 || go install mvdan.cc/gofumpt@v0.7.0
	gofumpt -w cmd internal tools
	golangci-lint --version | grep 1.64.5 || wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.5
	golangci-lint run --verbose  --max-issues-per-linter 0 --max-same-issues 0 --fix
