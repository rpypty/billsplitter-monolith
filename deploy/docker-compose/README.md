# Деплой через docker-compose.


### Поднять инфру и приложение:

1. Из текущей директории
```bash
    docker compose up --build -d
```

2. Из корневой директории через Makefile
```bash
make compose-up
make compose-down # для остановки
```


### Прогон миграций

1. Из текущей директории
```bash
    docker compose run --rm migrate
```

2. Из корневой директории через Makefile
```bash
    make compose-migrate
```


### Посмотреть логи приложения

```bash
docker compose -f deploy/docker-compose/docker-compose.yml logs -f
```