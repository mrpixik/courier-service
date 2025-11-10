up: # разворачивание и запуск
	docker compose up

down: # остановка всех контейнеров
	docker compose stop service-courier
	docker compose stop migrations
	docker compose stop postgres

copy_env_win: # копирование конфига (для Windows)
	copy /y .env.example .env

copy_env_lin: # копирование конфига (для Linux)
	cp -n .env.example .env
