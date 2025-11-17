PG_DSN="postgres://admin:admin@localhost:55433/bill_splitter?sslmode=disable&connect_timeout=5"
TEST_PKGS := $(shell go list ./... | grep -v mocks)

# Обновление зависимостей, проверка кода и обновление вендор папки
gop:
	go mod tidy && go mod vendor && go vet ./...

compose:
	docker-compose up --build -d

compose-down:
	docker-compose down

run:
	go run ./cmd/main.go

# Запускает миграции на postgres
migrate:
	goose -dir internal/db/migrations postgres $(PG_DSN) up

# Запускает миграции даже в случае непоследовательного применения файлов
migrate-force:
	goose -allow-missing -dir internal/db/migrations postgres $(PG_DSN) up

# Генерация swagger документации
swagger:
	swag init --parseDependency --parseInternal -g ./cmd/main.go

pre-commit:
	make gop
	make mocks
	make test
	make swagger

# Генерация моков
mocks:
	mockery --config .mockery.yaml

# Запуск тестов
test:
	go test -v -count=1 ./...

teste2e:
	go test -v -count=1 -tags=e2e ./...

# Посмотреть тестовое покрытие
test-cover:
	go test -covermode=atomic -coverpkg=./... -coverprofile=coverage.out $(TEST_PKGS)
	grep -v "/mocks/" coverage.out > coverage_filtred.out
	go tool cover -func=coverage_filtred.out
	rm coverage.out coverage_filtred.out
