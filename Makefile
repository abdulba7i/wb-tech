# Makefile for wb-tech

.PHONY: build app-dev app-stop app-dev-logs app-clean

build:
	docker compose build

app-dev:
	docker compose up --build

app-stop:
	docker compose down

app-dev-logs:
	docker compose logs -f

app-clean:
	docker compose down -v --rmi all --remove-orphans 