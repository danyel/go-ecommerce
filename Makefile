build:
	go build -o bin/main cmd/main.go

run:
	air

ui:
	cd gocommerce && npm run dev