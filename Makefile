run:
	@./script/bin.sh run
.PHONY: run

up:
	@docker compose up -d
.PHONY: up

down:
	@docker compose down
.PHONY: down

build:
	@docker compose build
.PHONY: build