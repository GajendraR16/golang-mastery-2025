SHELL := /bin/bash

DOCKER ?= docker
GO ?= go
COMPOSE = $(DOCKER) compose
TEST_DB_CONTAINER := task-db-test
TEST_DB_USER := postgres
TEST_DB_PASSWORD := postgres
TEST_DB_NAME := taskdb_test
TEST_DB_PORT := 5433

.PHONY: help up down restart logs ps build run local-db-up local-db-init local-db-down test test-cover test-handler test-db-up test-db-init test-db-down test-with-db clean

help:
	@echo "Available targets:"
	@echo "  make up            - Start API + Postgres with Docker Compose"
	@echo "  make down          - Stop Docker Compose services"
	@echo "  make restart       - Restart Docker Compose services"
	@echo "  make logs          - Tail Docker Compose logs"
	@echo "  make ps            - Show Docker Compose service status"
	@echo "  make build         - Build API image via Docker Compose"
	@echo "  make run           - Run API locally (go run main.go)"
	@echo "  make local-db-up   - Start local Postgres for app on localhost:5434 (container: taskdb)"
	@echo "  make local-db-init - Initialize local app DB schema"
	@echo "  make local-db-down - Remove local app DB container"
	@echo "  make test-db-up    - Start PostgreSQL test DB container on localhost:5433"
	@echo "  make test-db-init  - Initialize test DB schema"
	@echo "  make test-db-down  - Remove PostgreSQL test DB container"
	@echo "  make test          - Run all tests"
	@echo "  make test-cover    - Run tests with coverage"
	@echo "  make test-handler  - Run handler tests only"
	@echo "  make test-with-db  - Start test DB, init schema, run tests, cleanup"
	@echo "  make clean         - Stop compose services and remove test DB container"

up:
	$(COMPOSE) up --build

down:
	$(COMPOSE) down

restart:
	$(COMPOSE) down
	$(COMPOSE) up --build

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

build:
	$(COMPOSE) build

run:
	go run main.go

local-db-up:
	@$(DOCKER) rm -f taskdb >/dev/null 2>&1 || true
	$(DOCKER) run --name taskdb \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=taskdb \
		-p 5434:5432 \
		-d postgres:15

local-db-init:
	$(DOCKER) exec taskdb pg_isready -U postgres -d taskdb
	$(DOCKER) exec -i taskdb psql -U postgres -d taskdb < schema.sql

local-db-down:
	$(DOCKER) rm -f taskdb

test-db-up:
	@$(DOCKER) rm -f $(TEST_DB_CONTAINER) >/dev/null 2>&1 || true
	$(DOCKER) run --name $(TEST_DB_CONTAINER) \
		-e POSTGRES_USER=$(TEST_DB_USER) \
		-e POSTGRES_PASSWORD=$(TEST_DB_PASSWORD) \
		-e POSTGRES_DB=$(TEST_DB_NAME) \
		-p $(TEST_DB_PORT):5432 \
		-d postgres:15

test-db-init:
	@echo "waiting for $(TEST_DB_CONTAINER) to become ready..."
	@attempts=0; \
	until $(DOCKER) exec $(TEST_DB_CONTAINER) pg_isready -U $(TEST_DB_USER) -d $(TEST_DB_NAME) >/dev/null 2>&1; do \
		attempts=$$((attempts + 1)); \
		if [ $$attempts -ge 30 ]; then \
			echo "postgres did not become ready in time"; \
			$(DOCKER) logs $(TEST_DB_CONTAINER) | tail -n 50; \
			exit 1; \
		fi; \
		sleep 1; \
	done
	$(DOCKER) exec -i $(TEST_DB_CONTAINER) psql -U $(TEST_DB_USER) -d $(TEST_DB_NAME) < schema.sql

test-db-down:
	$(DOCKER) rm -f $(TEST_DB_CONTAINER)

test:
	$(GO) test ./...

test-cover:
	$(GO) test -cover ./...

test-handler:
	$(GO) test -v ./handler/

test-with-db:
	@set -e; \
	trap '$(MAKE) --no-print-directory test-db-down >/dev/null 2>&1 || true' EXIT; \
	$(MAKE) --no-print-directory test-db-up; \
	$(MAKE) --no-print-directory test-db-init; \
	$(MAKE) --no-print-directory test

clean:
	$(COMPOSE) down
	@docker rm -f $(TEST_DB_CONTAINER) >/dev/null 2>&1 || true
