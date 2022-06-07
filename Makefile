all: postgres

inmemory:
	@echo "start with inmemory storage"
	STORAGE_TYPE=inmemory docker-compose up -d --build

postgres:
	@echo "start with postgres storage"
	STORAGE_TYPE=postgres docker-compose --profile postgres up -d --build
