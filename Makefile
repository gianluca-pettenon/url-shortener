-include .env
export

COMPOSE := docker compose
IMAGE   := migrate/migrate:v4.19.1
URL     := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@127.0.0.1:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: migrate up down image shorten list load-test

migrate:
	$(COMPOSE) up -d --wait postgres
	docker run --rm --network container:postgres -v $(CURDIR)/migrations:/migrations \
		$(IMAGE) -path /migrations -database "$(URL)" up

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

image:
	$(COMPOSE) --profile cli build url

shorten:
	$(COMPOSE) --profile cli run --rm url ./shortener shorten

load-test:
	$(COMPOSE) --profile cli run --rm url ./shortener load-test

list:
	$(COMPOSE) --profile cli run --rm url ./shortener list
