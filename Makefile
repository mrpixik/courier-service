up_local: # запуск бд и миграций через докер. Само приложение запустится через go run
	docker-compose -f docker-compose.local.yaml up -d
	go run cmd/main.go --env=local

down_local: # остановка контейнеров. чтобы завершить работу сервиса необходимо еще отправить ctr+c в консоль
	docker compose -f docker-compose.local.yaml stop

up_prod: # запуск всего сервера (подгружается образ с моего dockerHub)
	docker-compose -f docker-compose.prod.yaml up -d

down_prod: # остановка всех контейнеров
	docker compose -f docker-compose.prod.yaml stop service-courier
	docker compose -f docker-compose.prod.yaml stop migrations
	docker compose -f docker-compose.prod.yaml stop postgres

run_tests:
	go test ./internal/http/server/handlers/courier
	go test ./internal/http/server/handlers/delivery
	go test ./internal/service/courier
	go test ./internal/service/delivery

run_tests_with_coverage:
	go test -cover ./internal/http/server/handlers/courier
	go test -cover ./internal/http/server/handlers/delivery
	go test -cover ./internal/service/courier
	go test -cover ./internal/service/delivery

run_test_integration:
	docker-compose -f docker-compose.tests.yaml up -d
	go test -v -tags=integration ./internal/repository/postgres/integration/courier
	go test -v -tags=integration ./internal/repository/postgres/integration/delivery
	docker-compose -f docker-compose.tests.yaml down -v

deploy_local: #для личного удобства
	docker build -t courier-service:latest -f ./deploy/docker/Dockerfile .

deploy_push_remote: #для личного удобства
	docker login
	docker tag courier-service:latest pixik/courier-service:latest
	docker push pixik/courier-service:latest
