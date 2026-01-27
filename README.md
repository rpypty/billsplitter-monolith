# Bill Splitter - бэк, монолит

## Локальный запуск

Перед запуском надо поднять инфру и накатить миграции (как это сделать в написано в соответствующих разделах инструкции)

```bash
make compose # поднимает инфру в докере 
make migrate # применяет миграции
make run     # запуска go app
```
Сервер будет доступен на порту из конфига (config.yml), по дефолту: http://localhost:5001

Сваггер доступен на http://localhost:5001/swagger/


## Конифг сервиса

Конфиг хранится в корне репозитория в файле config.yml

## Работа с БД
```bash
# Подключение к базе через CLI:
psql -h localhost -p 55433 -U admin -d bill_splitter
```

### Миграции

Миграции применяются с помощью goose - это go tool

#### 1. Установка мигратора
```bash
GOPROXY=direct go install github.com/pressly/goose/v3/cmd/goose@latest
```

#### 2. Создании файла миграции 

В момент прогона этой команды важно находиться в директории с миграциями: internal/db/migrations

```bash
goose -dir internal/db/migrations create <file_name> sql
```

#### 3. Применение миграции (up)
```bash
goose -dir internal/db/migrations postgres "postgres://admin:admin@localhost:55433/bill_splitter?sslmode=disable&connect_timeout=5" up
```
Через Makefile:
```bash
make migrate       # Прогоняет миграции
make migrate-force # Применяет миграции даже в случае непоследовательного применения файлов
```

## Инфраструктура, docker

Инфра поднимается через docker-compose.yml в корне репозитория

#### 4. Поднять инфру
```bash
docker-compose up --build -d
# или
make compose
```

## Swagger

#### 1. Установка

```bash
GOPROXY=direct go install github.com/swaggo/swag/cmd/swag@latest
swag --version
```

#### 2. Генерация

```bash
swag init --parseDependency --parseInternal -g ./cmd/main.go
```
Это создаст папку /docs с файлом docs.go.

Сваггер доступен на http://localhost:5001/swagger/

Для отправки запросов из сваггера надо:
- авторизоваться через роут /auth
- скопировать ID сессии
- проставить ID сессии в поле Authorize (сверху справа)

После этих шагов все остальные запросы будут авторизованы


## Генерация, линтинг, тесты
Линтинга пока нет

Для генерации доки, моков, пересборки зависимостей и прогона тестов предназначена команда:
```bash
make pre-commit
```
