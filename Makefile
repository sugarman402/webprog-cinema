.PHONY: help build up down logs clean ps db-reset restart rebuild

COMPOSE := docker compose

help:
	@echo "Cinema Ticket - Makefile Commands"
	@echo "==================================="
	@echo "make build          - Build all images"
	@echo "make up             - Start containers"
	@echo "make down           - Stop containers"
	@echo "make logs           - View logs"
	@echo "make clean          - Remove containers and data"
	@echo "make db-reset       - Reset database"
	@echo "make ps             - Show container status"
	@echo "make restart        - Restart containers"
	@echo "make rebuild        - Rebuild and start containers"

build:
	$(COMPOSE) build

up:
	$(COMPOSE) up -d
	@echo "Application starting: http://localhost"
	@echo "Backend API: http://localhost:8080"

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

clean:
	$(COMPOSE) down -v
	@echo "Containers and data removed"

ps:
	$(COMPOSE) ps

db-reset:
	$(COMPOSE) down -v
	$(COMPOSE) up -d postgres
	@echo "Database reset to initial state"

restart:
	$(COMPOSE) restart

rebuild:
	$(COMPOSE) build --no-cache
	$(COMPOSE) up -d
