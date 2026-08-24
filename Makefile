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
	$(COMPOSE) build url-shortener

shorten:
	$(COMPOSE) exec -it url-shortener ./shortener shorten

load-test:
	$(COMPOSE) exec -it url-shortener ./shortener load-test

list:
	$(COMPOSE) exec url-shortener ./shortener list
