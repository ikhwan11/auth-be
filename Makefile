APP_NAME := auth-be
MAIN_FILE := cmd/api/main.go

.PHONY: swagger run build tidy migrateup migratedown

swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g $(MAIN_FILE) -o docs

run:
	go run ./$(MAIN_FILE)

build:
	go build -o bin/$(APP_NAME) ./$(MAIN_FILE)

tidy:
	go mod tidy

migrateup:
	go run $(MAIN_FILE) migrate:up

migratedown:
	go run $(MAIN_FILE) migrate:down