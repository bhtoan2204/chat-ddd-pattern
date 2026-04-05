run:
	@echo "Formatting code..."
	@goimports -w .
	@script/bin.sh run
.PHONY: run

run-gateway:
	@cd gateway && make run
.PHONY: run-gateway

up:
	@docker compose up -d
.PHONY: up

down:
	@docker compose down
.PHONY: down

build:
	@docker compose build
.PHONY: build

generate:
	go run scaffold/main.go
.PHONY: generate