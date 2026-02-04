# Bill Splitter - бэк, монолит

## Локальный запуск

Перед запуском надо поднять инфру и накатить миграции (как это сделать в написано в соответствующих разделах инструкции)

```bash
make compose-infra-up # поднимает инфру в докере
make migrate # применяет миграции
make run     # запуска go app
```
Сервер будет доступен на порту из конфига (config.yml), по дефолту: http://localhost:5001

Сваггер доступен на http://localhost:5001/swagger/


## Конифг сервиса

Конфиг хранится в корне репозитория в файле config.yml.
Можно переопределять значения через `.env` (приоритет выше, чем у `config.yml`). Список переменных есть в `.env.example`.

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

#### Миграции в проде (deploy/docker-compose/docker-compose.yml)
Поднимаем всё через compose, затем накатываем миграции с хоста (goose установлен на сервере):
```bash
docker compose -f deploy/docker-compose/docker-compose.yml up -d
# если другие креды/порт — переопредели переменными
PG_HOST=127.0.0.1 PG_PORT=55433 PG_USER=admin PG_PASSWORD=admin PG_DB=bill_splitter make migrate
```
Либо можно прогнать миграции через отдельный сервис:
```bash
docker compose -f deploy/docker-compose/docker-compose.yml run --rm migrate
```

## Инфраструктура, docker

Compose-файлы:
- `deploy/docker-compose/docker-compose.yml` — полный запуск (приложение + БД)
- `docker-compose-infra.yml` — только инфраструктура (БД, для локального дебага)

#### 4. Поднять инфру
```bash
docker compose -f docker-compose-infra.yml up -d
# или
make compose-infra-up
```

#### 5. Поднять всё (прод-сценарий)
```bash
docker compose -f deploy/docker-compose/docker-compose.yml up --build -d
# или
make compose-up
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


## Как деплоить на VPS

При попытке поднимать инфру и приложение в докере возникает ошибка:
```
failed to solve: failed to fetch anonymous token: Get "https://auth.docker.io/token?scope=repository%3Alibrary%2Falpine%3Apull&service=registry.docker.io": net/http: TLS handshake timeout
```

Это скорее всего связано с региональными ограничениями - впс московский

Будет запускать приложение без докера просто в фоне:
```bash
go build -o billsplitter cmd/main.go # билдим бинарь
nohup ./billsplitter > /var/log/billsplitter.log 2>&1 & # запускаем в фоне
```

Чтобы дропнуть процесс:
```bash
pkill -9 -f "./billsplitter" # остановить процесс
```
