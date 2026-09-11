# Архитектура проекта

Короткая справка по текущему устройству `/home/tagir/bmstu/rsoi`.

## Runtime Flow

```text
cmd/service
  -> app.New
       init logging/tracing/metrics
       init db/rest/core metrics
       create app.Manager
  -> app.Manager.Start
       db.Manager.Prepare()
       db.Exports()
       persons.NewPersonUseCase(dbExports.PersonRepository)
       status.NewStatusUseCase(dbExports.StatusProvider)
       rest.Manager.Prepare(personUseCase, statusUseCase)
       db.Manager.Run()
       rest.Manager.Run()
```

Shutdown идет в обратном порядке через `app.Close`: REST, DB, tracing provider, logger.

## Packages

```text
internal/app       application wiring and lifecycle
internal/config    config structs and loader
internal/core   use cases and their metrics
internal/db        PostgreSQL/Ent client, repositories, DB health provider
internal/models    cross-layer interfaces and entities
internal/rest      REST service wiring, Huma routes, HTTP servers
internal/observability
tests/integration database-level tests
tests/e2e         black-box HTTP tests
```

`internal/core` contains business use cases. `internal/rest/restsvc/resources/persons` is only the external HTTP resource shape for `/api/v1/persons`.

## Dependency Wiring

There is no dependency container. `app.Manager` wires dependencies explicitly:

```text
db exports repository/status provider
core constructors build use cases
rest receives use cases in Prepare
```

This keeps ownership visible and avoids a service locator.

## Persistence

DB runtime uses Ent over `database/sql` with pgx:

```text
internal/db/ent/
internal/db/repos/
internal/db/resources/status_provider.go
migrations/
```

SQL migrations are applied with Atlas.

## HTTP API

Base prefix:

```text
/api/v1
```

Endpoints:

```text
GET    /persons
GET    /persons/{id}
POST   /persons
PATCH  /persons/{id}
DELETE /persons/{id}

GET /status
GET /healthz
GET /readyz
GET /version
```

## Metrics

Global metrics use the application/service prefix. Module metrics include:

```text
db_connections_active
db_connections_created_total
db_connections_destroyed_total
db_connections_alive
rest_service_starts_total
core_persons_*_request_total
core_persons_*_request_failed_total
status_get_status_request_total
status_healthz_request_total
status_readyz_request_total
```

## Useful Commands

```bash
make generate
make fmt
make lint
make test
make integration
make e2e
```

Sandbox-friendly checks:

```bash
GOCACHE=/tmp/rsoi-go-build-cache go test ./...
GOCACHE=/tmp/rsoi-go-build-cache GOLANGCI_LINT_CACHE=/tmp/rsoi-golangci-lint-cache make lint
```
