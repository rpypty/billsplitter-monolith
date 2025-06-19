# Bill Splitter - бэк, монолит

## Локальный запуск

```bash
make compose # поднимает инфру в докере 
make run     # запуска go app
```
Сервер будет доступен на порту из конфига (config.yml), по дефолту: http://localhost:5001

Сваггер доступен на http://localhost:5001/swagger/


## Конифг

Конфиг храинтся в файле config.yml

## Миграции

#### 1. Установка мигратора
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

#### 2. Создании файла миграции 
```
goose -dir internal/db/migrations create <file_name> sql
```

#### 3. Применение миграции (up)
```
goose -dir internal/db/migrations postgres "postgres://seller:seller@localhost:55433/bill_splitter?sslmode=disable&connect_timeout=5" up
```

## Докер

#### 4. Поднять инфру
```bash
docker-compose up --build -d
```

## Swagger

#### 1. Установка

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag --version
```

#### 2. Генерация

```bash
swag init --parseDependency --parseInternal -g ./cmd/main.go
```
Это создаст папку /docs с файлом docs.go.

