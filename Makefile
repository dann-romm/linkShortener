all: postgres

inmemory:
	@echo "start with inmemory storage"
	docker-compose up -d --build -e STORAGE_TYPE=inmemory

postgres:
	@echo "start with postgres storage"
	docker-compose up -d --build --profile postgres -e STORAGE_TYPE=postgres
