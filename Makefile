PG_DSN="postgres://seller:seller@localhost:55433/bill_splitter?sslmode=disable&connect_timeout=5"

gop:
	go mod tidy && go mod vendor && go vet ./...

compose:
	docker-compose up --build -d

run:
	go run ./cmd/main.go

migrate:
	goose -dir internal/db/migrations postgres $(PG_DSN) up

swagger:
	swag init --parseDependency --parseInternal -g ./cmd/main.go

