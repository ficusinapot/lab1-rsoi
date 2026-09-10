# Продолжение по аналогии с bsdn-core-v2

Этот файл нужен как контекст для новых чатов. Проект в `/home/tagir/bmstu/rsoi` приводится по архитектурной аналогии с `/home/tagir/bsdn/bsdn-core-v2`, но без закрытых common-пакетов из BSDN. Вместо них сделаны локальные минимальные аналоги.

## Эталон

Эталонный проект:

```text
/home/tagir/bsdn/bsdn-core-v2
```

С него берется форма модулей:

```text
internal/<module>/
  config.go
  errors.go
  imports.go
  exports.go
  components.go
  manager.go
  metrics.go
  metrics/metrics.go
```

И общий lifecycle:

```text
Prepare -> Run -> StopCommunications -> Shutdown -> Dispose
```

## Что уже сделано

Основная цепочка managers:

```text
app.Manager
  -> db.Manager
       exports PersonRepository
       exports StatusProvider
  -> domain.Manager
       imports PersonRepository
       exports PersonUseCase
  -> status.Manager
       imports StatusProvider
       exports StatusUseCase
  -> rest.Manager
       imports PersonUseCase
       imports StatusUseCase
       starts REST service
```

Добавлен общий lifecycle-пакет:

```text
internal/manager/
  errors.go
  state.go
  manager.go
  manager_test.go
```

Он содержит:

- `Manager` interface
- `Inner` interface
- `GenericManager`
- `State`
- `ValidateTransition`
- общие ошибки `prepare`, `run`, `stop_communications`, `shutdown`, `invalid_state`

Модули `db`, `domain`, `status`, `rest` переведены на `manager.GenericManager`.

Добавлен локальный dependency container:

```text
internal/deps/
  container.go
  container_test.go
```

Managers больше не получают ссылки на соседние managers в constructors. `app.Manager` накапливает typed exports через `deps.Container` и передает их в `PrepareWithDeps`:

```text
db.Prepare(empty) -> db exports
domain.Prepare(db exports) -> domain exports
status.Prepare(db exports) -> status exports
rest.Prepare(domain + status exports)
```

Подсистема `persons` переименована в `domain`, чтобы bounded context не был привязан к одной сущности. REST resource `/persons` сохранен как внешний API-контракт.

Внутри domain выделен submodule для текущей сущности:

```text
internal/domain/persons/
  person.go
  person_create.go
  person_get.go
  person_list.go
  person_update.go
  person_delete.go
  metrics/metrics.go
```

Build metadata для `/api/v1/version` заполняется через `Makefile` и `-ldflags`:

```text
build_time
branch
commit
```

REST-модули имеют README и OpenAPI tags. В `internal/rest/restsvc/routes.go` добавлен central error hook для сообщений формата `Title: details`.

REST resources выровнены на одном уровне иерархии:

```text
internal/rest/restsvc/resources/
  persons/
  status/
  version/
```

E2E-тесты разделены по пакетам:

```text
tests/e2e/persons/
tests/e2e/status/
tests/e2e/version/
tests/e2e/metrics/
tests/e2e/openapi/
tests/e2e/infrastructure/
```

Общая инфраструктура для integration и e2e вынесена отдельно:

```text
tests/infrastructure/
```

DB runtime переписан с GORM на Ent:

```text
internal/db/ent/
  schema/person.go
  generate.go
```

`db.Open` открывает `database/sql` через pgx, создает Ent client поверх PostgreSQL dialect и сохраняет `*sql.DB` для метрик. SQL migrations остаются в `migrations/` и применяются через Atlas.

## Реализованные REST endpoints

Base prefix:

```text
/api/v1
```

Persons CRUD:

```text
GET    /api/v1/persons
GET    /api/v1/persons/{id}
POST   /api/v1/persons
PATCH  /api/v1/persons/{id}
DELETE /api/v1/persons/{id}
```

Status:

```text
GET /api/v1/status
GET /api/v1/healthz
GET /api/v1/readyz
```

Version:

```text
GET /api/v1/version
```

OpenAPI:

```text
GET /api/v1/openapi.yaml
GET /api/v1/openapi.json
GET /api/v1/docs
GET /api/v1/swagger
```

`/docs` и `/swagger` включаются через:

```yaml
rest:
  service:
    openapi:
      ui: true
```

## Реализованные метрики

Общие:

```text
service_app_info
service_rest_api_http_requests_total
service_rest_api_http_request_duration_seconds
service_db_connections_open
service_db_connections_in_use
service_db_connections_idle
```

Модульные:

```text
db_connections_active
db_connections_created_total
db_connections_destroyed_total
db_connections_alive
rest_service_starts_total
domain_persons_create_request_total
domain_persons_create_request_failed_total
domain_persons_get_request_total
domain_persons_get_request_failed_total
domain_persons_list_request_total
domain_persons_list_request_failed_total
domain_persons_update_request_total
domain_persons_update_request_failed_total
domain_persons_delete_request_total
domain_persons_delete_request_failed_total
status_get_status_request_total
status_healthz_request_total
status_readyz_request_total
```

## Конфиг

Текущая структура `configs/config.yaml`:

```yaml
rest:
  service:
    addr: ":8080"
    timeout:
      shutdown: "30s"
      read: "0s"
      write: "0s"
      idle: "0s"
    openapi:
      ui: true
  metrics:
    enabled: true
    addr: "0.0.0.0:11190"

database:
  dsn: "host=localhost user=program password=test dbname=persons port=5432 sslmode=disable"
  max_open_conns: 10
  max_idle_conns: 5
  conn_max_lifetime: "1h"
  health_check:
    check_period: "10s"
```

## Тесты, которые должны проходить

После каждого шага запускать:

```bash
go test ./...
go test -count=1 -tags='integration e2e' ./...
go mod tidy -diff
make lint
```

Integration/e2e используют testcontainers и PostgreSQL.

Для локальной разработки можно поднять PostgreSQL 18 через Docker Compose:

```bash
cd ci
docker compose up -d
```

Compose создает БД `persons` и пользователя `program` с паролем `test`.

## Что делать дальше

Следующий крупный шаг явно не зафиксирован. Перед продолжением стоит выбрать новый архитектурный фокус: конфигурация, observability, REST error contracts или разбиение db слоя.

## Важные ограничения

- Не подтягивать закрытые BSDN common-пакеты.
- Не откатывать существующие незакоммиченные изменения.
- Если меняется REST/OpenAPI DTO, следить за конфликтами имен типов в Huma. Например, `status.Response` конфликтовал с `persons.Response`, поэтому status DTO названы `StatusResponse`, `GetStatusResponse`.
