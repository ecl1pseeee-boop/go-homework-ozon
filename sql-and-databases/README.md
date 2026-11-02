### ДЗ4

1. В директории docker-entrypoint-initdb.d лежат sql для инициализации БД. Они запускаются при старте БД.
2. Запуск контейнера с PostgreSQL
```shell

docker run --name pg-hw \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=homework \
  -v ./init-scripts:/docker-entrypoint-initdb.d \
  -p 5432:5432 \
  -d postgres:latest
 
 ```
3. Подключение к PostgreSQL для удобства написания, выполнения и тестирования SQL запросов я делал внутри Goland. 
Внутри GoLand есть встроенная вкладка Database. Достаточно подключиться к localhost:5432 (пользователь postgres, пароль postgres, БД homework).
