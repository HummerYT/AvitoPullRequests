.PHONY: up down test load-test clean build

up:
	docker build -t avitopullrequest .
	docker compose up -d

down:
	docker compose down

rebuild: down up

load-test:
	chmod +x load_test.sh
	./load_test.sh

build:
	go build -o bin/app cmd/main.go

run: build
	./bin/app

logs:
	docker compose logs -f app

status:
	docker compose ps