# Обёртка над пакетом "log/slog"

Предназначена для локальной разработки

Предоставляет удобочитаемый формат отображения логов в консоль с цветовым  дифференцированием. Для удобства чтения используется отступ каждого уровня вложенности "\t"

Минимальный отображаемы уровень логов "DEBUG"

## Установка
```sh
go get github.com/rautaruukkipalich/prettyslog/@latest
```

## Инициализвация логгера
```
log := prettyslog.NewPrettyLogger()
```

## Пример отображения логов
![alt test1](https://github.com/rautaruukkipalich/prettyslog/blob/main/img/1.PNG?raw=true)