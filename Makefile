build:
	go build -o bin/main cmd/main.go

run:
	air

ui:
	cd gocommerce && npm run dev

migrate:
	goose -dir ./migrations postgres "$$DATABASE_URL" up