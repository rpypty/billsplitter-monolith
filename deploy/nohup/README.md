### Деплой через создание процесса в бэкграунде


Запуск (из корня репозитория)
```bash
nohup ./billsplitter > /var/log/billsplitter.log 2>&1 &
```

Остановка
```bash
pkill -9 -f "./billsplitter"
```
