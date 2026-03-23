# Team-Task-Manager

Task manager service for team work with change history

## Requirements
GNU Make 3.81, g1.26.1, Docker 28.1.1 (Docker Desktop)

### In `secrets/` directory placed default secrets. At first you can change it.

## Run local
1. Prepare local database and cache
```shell
make start-local-database
make migrations-up
make start-local-redis
```

2. Run local service
```shell
make run-local
```

## Run Docker Compose
1. Run servce
```shell
make service
```

## Explore via OpenAPI
When service is running the OpenAPI page with all implemented endpoints is available at [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## Makefile targets
- `make deps` download all required libraries
- `make errcheck` scan code for some errors
- `make linter` run linter
- `make generate-sqlc` generates Go code by sql queries. sqlc config is in `sqlc.yaml` file
- `make generate-mocks` generate mocks for interfaces specified in `//go:generate ...` comments in any Go file
- `make generate-swag` generates OpenAPI docs from comments of endpoint functions
- `make migrations-up` applies all migrations
- `make migrations-down` rollbacks of last migration
- `make migrations-down-all` rollback to 0 migration version
- `make migrations-status` shows migrations status
- `make start-local-database` runs local database
- `make stop-local-database` stops local database
- `make clean-local-database` cleans local database
- `make start-local-redis` runs local redis
- `make stop-local-redis` stops local redis
- `make clean-local-redis` cleans local redis
- `make run-local` compiles and runs service binary
- `make run-local-fast` runs service via `go run` command
- `make service` runs docker compose service
- `make service-stop` stops docker compose service
- `make service-rebuild` rebuild docker compose image
- `make coverage-info` shows unit testing coverage in terminal
- `make coverage-html` shows unit testing coverage in browser
- `make unit-test` runs unit tests
- `make integration-test` runs integration tests
- `make clean` removes compiled binaries
