build:
	go build -o bin/main cmd/main.go

env_up:
	docker compose up -d

env_down:
	docker compose down

run:
	air

ui:
	cd gocommerce && npm run dev

migrate:
	goose up

full:
	go build -v ./.. & go test -v ./...

tools:
	go install github.com/pressly/goose/v3/cmd/goose@latest && \
	go install github.com/air-verse/air@latest

integration_test:
	go test github.com/danyel/ecommerce/test/integration

mock_tests:
	go test github.com/danyel/ecommerce/test/mock