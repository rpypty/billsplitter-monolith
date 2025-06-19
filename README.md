# Bill Splitter - бэк, монолит

## Запуск
```bash
go build -o bin/billsplitter ./cmd
go run ./cmd/main.go
```

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