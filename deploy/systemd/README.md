# Как задеплоить бэк на VPS

Изначально предполагается что репозиторий уже лежит на впс и бинарник компилится 

## Ручная настройка и деплой

### 1. Подготовка машины, содзание пользователя и служебных папок (делается один раз)

```bash
sudo adduser --system --group --no-create-home billsplitter   # отдельный пользователь без логина (безопаснее)
sudo mkdir -p /opt/billsplitter                               # каталог для приложения

```
### 2. Билдим приложение и переносим в служебную папку (находимся в директории репозитория)

```bash
go build -o billsplitter cmd/main.go
sudo mv ./billsplitter /opt/billsplitter/billsplitter         # перенести бинарник в /opt
sudo chown -R billsplitter:billsplitter /opt/billsplitter     # отдать права приложению
sudo chmod 755 /opt/billsplitter/billsplitter                 # сделать исполняемым
```

### 3*. Копируем конфиг (при необходимости)

```bash
cp ./config.yml /opt/billsplitter/config.yml                   # копируем конфиг
sudo chown billsplitter:billsplitter /opt/billsplitter/config.yml  # владелец сервисный юзер
sudo chmod 600 /opt/billsplitter/config.yml                        # доступ только владельцу
```


### 4. Настраиваем systemd unit (один раз)

#### 4.2 Создаем службу

```bash
sudo nano /etc/systemd/system/billsplitter.service         # создать unit (описание сервиса)
```

Содержимое сервиса:
```md
[Unit]
Description=bill-splitter backend
After=network-online.target                               
Wants=network-online.target                                

[Service]
User=billsplitter                                         # запускать не под root
Group=billsplitter
WorkingDirectory=/opt/billsplitter                         # рабочая директория (если нужны файлы рядом)

ExecStart=/opt/billsplitter/billsplitter                   # команда запуска (добавишь флаги если нужно)
Restart=always                                             # перезапуск при падении
RestartSec=2                                               # пауза между перезапусками

# Логи по умолчанию пойдут в journald, читать через journalctl
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target 
```

#### 4.2 Применить unit, запустить, включить автозапуск

```bash
sudo systemctl daemon-reload                                # перечитать новые unit-файлы
sudo systemctl enable billsplitter                          # включить автозапуск при ребуте
sudo systemctl start billsplitter                           # старт сервиса прямо сейчас
sudo systemctl status billsplitter --no-pager               # посмотреть статус без "пейджера"
```


### 5. Полезные команды

Посмотреть логи
```bash
journalctl -u billsplitter -f
journalctl -u billsplitter -n 200 --no-pager
sudo systemctl restart billsplitter                         # перезапуск
sudo systemctl stop billsplitter                            # остановка
sudo systemctl disable --now billsplitter                   # выключить автозапуск и остановить
```


## Депой из скрипта 

В репозитории в папке /deploy/systemd есть скрипт deploy.sh

```bash
sudo chmod +x ./deploy/systemd/deploy.sh    # делаем файл исполняемым
sudo ./deploy/systemd/deploy.sh             # запускаем скрипт деплоя
```
