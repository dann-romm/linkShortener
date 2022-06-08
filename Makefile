all: postgres

inmemory:
	@echo "start with inmemory storage"
	@mkdir -p ./data/inmemory
	STORAGE_TYPE=inmemory docker-compose up --build

postgres:
	@echo "start with postgres storage"
	STORAGE_TYPE=postgres docker-compose --profile postgres up --build

stop:
	@echo "stop all containers"
	docker-compose stop
