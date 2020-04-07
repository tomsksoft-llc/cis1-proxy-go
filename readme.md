## Сборка

```
$ cd {$project_path}/cmd/cis1-proxy-go
$ go install
```

Исполняемый файл будет находиться в ```$GOBIN```.

## Запуск

Вывод информации о запуске:

```
$ ./cis1-proxy-go [-h | --help]
```

Режим запуска:

```
$ ./cis1-proxy-go -a=[proxy_host] -p=[proxy_port] -c=[config_file_path] 
```

## Конфигурация маршрутов

Формат конфигурации:

```
[
	{
		"Location": "/path",
		"Pass": {
			"Host": "www.example.com",
			"Port": "80",
			"Path": "/"
		}
	},
	{
		"Location": "/another/path",
		"Pass": {
			"Host": "127.0.0.1",
			"Port": "8080",
			"Path": "/test/"
		}
	}
]
```

При такой конфигурации запросы будут проксироваться на:

```
http://proxy_host:proxy_port/path/{$whatever} -> http://www.example.com:80/{$whatever}
http://proxy_host:proxy_port/another/path/{$whatever} -> http://127.0.0.1:8080/test/{$whatever}
```

Пример конфига: ```cmd/cis1-proxy-go/config.json```.

