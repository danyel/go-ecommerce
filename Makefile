du:
	docker compose up -d

docker_up:
	docker compose up -d

dd:
	docker compose down

docker_down:
	docker compose down

bbe:
	go build -o bin/main cmd/main.go

build_backend:
	go build -o bin/main cmd/main.go

be:
	air

run_backend:
	air

bui:
	cd gocommerce && npm install --force

build_ui:
	make bui

ui:
	cd gocommerce && npm run dev

run_ui:
	make ui

dm:
	goose up

database_migration:
	make dm

ft:
	go build -v ./... && go test ./test/integration && go test ./test/mock

full_tests:
	make ft

inst:
	go install github.com/pressly/goose/v3/cmd/goose@latest && go install github.com/air-verse/air@latest

install_tools:
	make inst

it:
	go test github.com/danyel/ecommerce/test/integration

integration_tests:
	make it

mt:
	go test github.com/danyel/ecommerce/test/mock

mock_tests:
	make mt