Собрать контейнер
``` bash
docker buildx build -t andreybackend .
```
Запустить контейнер
``` bash
docker run -p 8080:8080 andreybackend
```

user = postgress
password = pgpass