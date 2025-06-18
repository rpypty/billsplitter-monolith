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
goose -dir db/migrations postgres "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up
```

