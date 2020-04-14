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
$ ./cis1-proxy-go -a=[proxy_host] -p=[proxy_port] -c=[config_file_path] -d=[cis_base_dir]
```

## Конфигурация

### 1. Маршрутизация

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

### 2. Запуск ядра напрямую

Формат конфигурации:

```
[
	{
		"Job": "/some-job",
		"Run": {
			"Project": "project_name",
			"Job": "job_name",
			"Args": ["arg1", "arg2", "arg3"]
		}
	},
	{
		"Job": "/another-job",
		"Run": {
			"Project": "another_project",
			"Job": "another_job",
			"Args": ["arg1", "arg2", "arg3"]
		}
	}
]
```

При такой конфигурации запросы будут запускать:

```
http://proxy_host:proxy_port/some-job -> ${cis_base_dir}/core/startjob project_name/job_name
http://proxy_host:proxy_port/another-job -> ${cis_base_dir}/core/startjob another_project/another_job --params arg1 arg2 arg3
```

Пример конфига: ```configs/config.json```.

