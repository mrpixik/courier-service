up_local: # запуск бд и миграций через докер. Само приложение запустится через go run
	docker-compose -f docker-compose.local.yaml up -d
	go run cmd/main.go --env=local

down_local: # остановка контейнеров. чтобы завершить работу сервиса необходимо еще отправить ctr+c в консоль
	docker compose -f docker-compose.local.yaml stop

up_prod: # запуск всего сервера (подгружается образ с моего dockerHub)
	docker-compose -f docker-compose.prod.yaml up -d

down_prod: # остановка всех контейнеров
	docker compose stop service-courier
	docker compose stop migrations
	docker compose stop postgres

deploy_local: #для личного удобства
	docker build -t courier-service:latest -f ./deploy/docker/Dockerfile .

deploy_push_remote: #для личного удобства
	docker login
	docker tag courier-service:latest pixik/courier-service:latest
	docker push pixik/courier-service:latest
