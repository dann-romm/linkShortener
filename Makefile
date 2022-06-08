all: postgres

build:
	@echo "building..."
	docker-compose build

inmemory:
	@echo "start with inmemory storage"
	@mkdir -p ./data/inmemory
	STORAGE_TYPE=inmemory docker-compose up -d

postgres:
	@echo "start with postgres storage"
	STORAGE_TYPE=postgres docker-compose --profile postgres up -d

stop:
	@echo "stop all containers"
	docker-compose stop
