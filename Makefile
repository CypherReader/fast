.PHONY: run build test migrate-up migrate-create

run:
	go run cmd/server/main.go

build:
	go build -o server.exe cmd/server/main.go

test:
	go test ./...

migrate-create:
	@if [ -z "$(NAME)" ]; then echo "NAME is required"; exit 1; fi
	@echo "Creating migration: migrations/$(shell date +%Y%m%d%H%M%S)_$(NAME).sql"
	@touch migrations/$(shell date +%Y%m%d%H%M%S)_$(NAME).sql
