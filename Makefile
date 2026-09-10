.PHONY: generate build migrate-up migrate-hash migrate-status unit-test integration-test e2e-test test-all lint fmt fmt-check tidy vuln docker-build test check clean

LINT_CONFIG   ?= ci/code-check/lint.yaml
FMT_CONFIG    ?= ci/code-check/fmt.yaml
GOLANGCI_LINT ?= $(if $(wildcard $(HOME)/go/bin/golangci-lint),$(HOME)/go/bin/golangci-lint,golangci-lint)
GO_VERSION    ?= $(shell go env GOVERSION)
GOVULNCHECK   ?= GOTOOLCHAIN=$(GO_VERSION) go run golang.org/x/vuln/cmd/govulncheck@latest
ATLAS         ?= atlas
MIGRATIONS_DIR ?= migrations
DB_DSN        ?= postgres://program:test@localhost:5432/persons?sslmode=disable
DOCKER_IMAGE  ?= rsoi-person-service
VERSION_PKG   := github.com/ficusinapot/ds/cmd/service/version
BUILD_TIME    ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BRANCH        ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)
COMMIT        ?= $(shell git rev-parse --verify HEAD 2>/dev/null || echo unknown)
LDFLAGS       := -X '$(VERSION_PKG).buildTime=$(BUILD_TIME)' -X '$(VERSION_PKG).branch=$(BRANCH)' -X '$(VERSION_PKG).commit=$(COMMIT)'

generate:
	go generate -v ./...

migrate-up:
	$(ATLAS) migrate apply --dir "file://$(MIGRATIONS_DIR)" --url "$(DB_DSN)"

migrate-hash:
	$(ATLAS) migrate hash --dir "file://$(MIGRATIONS_DIR)"

migrate-status:
	$(ATLAS) migrate status --dir "file://$(MIGRATIONS_DIR)" --url "$(DB_DSN)"

build:
	mkdir -p bin
	go build -o bin -v -ldflags "$(LDFLAGS)" ./...

unit-test:
	go test -race -timeout 30s ./internal/... ./cmd/...

integration-test:
	go test -count=1 -race -tags=integration -timeout 2m ./tests/integration/...

e2e-test:
	go test -count=1 -race -tags=e2e -timeout 2m ./tests/e2e/...

test-all: unit-test integration-test e2e-test

lint: build
	$(GOLANGCI_LINT) run -c $(LINT_CONFIG) ./...

fmt:
	$(GOLANGCI_LINT) fmt -c $(FMT_CONFIG)

fmt-check:
	$(GOLANGCI_LINT) fmt --diff -c $(FMT_CONFIG) ./...

tidy:
	go mod tidy --diff

vuln: build
	$(GOVULNCHECK) -show verbose ./...

docker-build:
	docker build -t $(DOCKER_IMAGE) .

check: tidy fmt-check lint vuln unit-test integration-test e2e-test build docker-build

clean:
	rm -rf bin
	go clean -testcache
