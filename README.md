# Обёртка над пакетом "log/slog"

Предназначена для локальной разработки

Предоставляет удобочитаемый формат отображения логов в консоль с цветовым  дифференцированием. Для удобства чтения используется указанный отступ каждого уровня вложенности

Минимальный отображаемы уровень логов "DEBUG"

## Установка
```sh
go get github.com/rautaruukkipalich/prettyslog@latest
```

## Инициализвация логгера

#### Выбрать нужный отступ
```
indent := "\t"
```

#### Инициализировать логгер передав в него "indent"
```
log := prettyslog.NewPrettyLogger(indent)
```
#### Или сразу передав его в функцию без объявления дополнительной переменной
```
log := prettyslog.NewPrettyLogger("\t")
```

## Пример отображения логов
![alt img1](https://github.com/rautaruukkipalich/prettyslog/blob/main/img/11.PNG?raw=true)
![alt img2](https://github.com/rautaruukkipalich/prettyslog/blob/main/img/22.PNG?raw=true)
