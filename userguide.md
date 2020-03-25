## 1. Маршрутизация

Роутинг подобный nginx, т.е. где-то в конфигурации описаны или при помощи regex или структурно пути и соотвествующие им сервисы.

Пример конфигурации:

```
proxy {
  path_a: http://128.127.126.126:8080/
  path_b: http://127.0.0.1:8123/test/
}
```

При такой конфигурации запросы будут проксироваться на:

```
http://your_host:yout_port/proxy/path_a/{$whatever} -> http://128.127.126.126:8080/{$whatever}
http://your_host:yout_port/proxy/path_b/{$whatever} -> http://127.0.0.1:8123/test/{$whatever}
```

Формат конфиграции на усмотрение разработчика.

## 2. Запуск ядра напрямую

* авторизация по токену: gitlab, github, raw. (аналогично [webui](https://github.com/tomsksoft-llc/cis1-webui-native-srv-cpp/wiki/Webhooks-tutorial), только с ручной конфигурацией)
* фильтрация (по маске или regex)
* запуск локально https://github.com/tomsksoft-llc/cis1-docs/blob/master/overview.md
* запуск по ssh (аналогично)
* передавать параметры, заданные в path