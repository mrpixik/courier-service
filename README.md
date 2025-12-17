# Сервис курьеров (Avito course)


TODO:
1. Использовать Strings.Builder при динамическом создании запроса

Что сделано (дз7):
1. в ./proto размещен .proto файл с описанием методов и структур для gRPC взаимодействия с service-order
2. в ./api/order помещены сгенерированные файлы для gRPC клиента и сервера
3. в ./internal/gateway/order добавлен адаптер для взаимодействия с gRPC клиентом service-order
4. в ./internal/worker/order реализован воркер, использующий адаптер (3) для RPC в service-order по тикеру
5. в конфиг добавлена настройка DSN подключения к gRPC серверу service-order

Что сделано (дз8):
1. новый service в internal/service/queues/order для обработки сообщений, получаемых через брокер (кафку)
2. сервис реализован с помощью 2х ПП: Fabric + Strategy
3. новый handler в internal/handler/queues/order для обработки сообщений
4. написан (скопирован с воркшопа) клиент для кафки в infrastructure/kafka/client
5. новый worker для фоновой обработки сообщений internal/worker/queues/kafka
6. все настраивается через конфиг, worker поддерживает graceful shutdown

Что сделано (дз9):
1. новый http observer, собирающий метрики через prometheus в internal/observability/metrics/prometheus
2. подключается к middleware в internal/handler/http/middleware
3. контейнер с prometheus и grafana поднимается в докер компоуз в infrastructure/monitoring

чтобы запустить сервис со всеми зависимостями (нужно, чтобы в корневой дирректории проекта лежал service-order)
его попросили не пушить в гитхаб, поэтому нужно скопировать его туда отдельно.
```bash
make up_prod
```
чтобы запустить unit тесты
```bash
make run_tests
```
либо с выводом покрытия
```bash
make run_tests_with_coverage
```

для интеграционных тестов
```bash
make run_test_integration
```