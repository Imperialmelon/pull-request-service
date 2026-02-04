# Pull Request Service

````markdown
## Запуск проекта

Для поднятия всех сервисов выполните команду:

```bash
docker-compose up -d
````

## Мониторинг

### Prometheus

Сбор метрик осуществляется с помощью Prometheus, доступного по адресу:

```
http://localhost:9090
```

### Grafana

Для визуализации метрик используется Grafana:

```
http://localhost:3000
```

Чтобы подключить Prometheus в Grafana как источник данных, используйте URL:

```
http://prometheus:9090
```

```
```
