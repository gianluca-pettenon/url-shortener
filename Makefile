-include .env
export

COMPOSE := docker compose --env-file .env
IMAGE   := migrate/migrate:v4.18.3
URL     := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@127.0.0.1:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: migrate migrate-down up down

migrate:
	$(COMPOSE) up -d --wait postgres
	docker run --rm --network container:postgres -v $(CURDIR)/migrations:/migrations \
		$(IMAGE) -path /migrations -database "$(URL)" up

migrate-down:
	$(COMPOSE) up -d --wait postgres
	docker run --rm --network container:postgres -v $(CURDIR)/migrations:/migrations \
		$(IMAGE) -path /migrations -database "$(URL)" down 1

up:
	$(COMPOSE) up -d --no-deps redis dns url

down:
	$(COMPOSE) down
