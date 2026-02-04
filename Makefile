PG_USER ?= admin
PG_PASSWORD ?= admin
PG_HOST ?= localhost
PG_PORT ?= 55433
PG_DB ?= bill_splitter
PG_SSLMODE ?= disable

PG_DSN := "postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=$(PG_SSLMODE)&connect_timeout=5"
TEST_PKGS := $(shell go list ./... | grep -v mocks)



# 1. Golang

# Обновление зависимостей, проверка кода и обновление вендор папки
run:
	go run ./cmd/main.go

gop:
	go mod tidy && go mod vendor && go vet ./...



# 2. Docker

# поднимает инфру + бэк
compose-up:
	docker compose -f deploy/docker-compose/docker-compose.yml up --build -d

# прогоняет миграции для постгри внутри контейнера (зависимость goose поднимается в докере)
compose-migrate:
	docker compose -f deploy/docker-compose/docker-compose.yml run --rm migrate

# посмотреть логи бэка поднятого в контейнере
compose-logs:
	docker compose -f deploy/docker-compose/docker-compose.yml logs -f app

# снести инфру + бэк
compose-down:
	docker compose -f deploy/docker-compose/docker-compose.yml down


# поднимает только инфру, для дева/дебага бэка
compose-infra-up:
	docker compose -f docker-compose-infra.yml up -d

compose-infra-down:
	docker compose -f docker-compose-infra.yml down



# 3. Migrations 

# Запускает миграции на postgres
migrate:
	goose -dir internal/db/migrations postgres $(PG_DSN) up

# Запускает миграции даже в случае непоследовательного применения файлов
migrate-force:
	goose -allow-missing -dir internal/db/migrations postgres $(PG_DSN) up



# 4. Утилиты

# Генерация моков
mocks:
	mockery --config .mockery.yaml

# Генерация swagger документации
swagger:
	swag init --parseDependency --parseInternal -g ./cmd/main.go

pre-commit:
	make gop
	make mocks
	make test
	make swagger

# Запуск тестов
test:
	go test -v -count=1 ./...

test-e2e:
	go test -v -count=1 -tags=e2e ./...

# Посмотреть тестовое покрытие
test-cover:
	go test -covermode=atomic -coverpkg=./... -coverprofile=coverage.out $(TEST_PKGS)
	grep -v "/mocks/" coverage.out > coverage_filtred.out
	go tool cover -func=coverage_filtred.out
	rm coverage.out coverage_filtred.out
