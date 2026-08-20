include .env
export

db-up:
	docker-compose up -d postgres

db-down:
	docker-compose down

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

migrate-down-all:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

migrate-create:
ifndef name
	$(error name is not set. Usage: make migrate-create name=add_status_column)
endif
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migrate-version:
	docker exec -it url-shortener-postgres-1 psql -U postgres -d shortener -c "SELECT * FROM schema_migrations;"

run: migrate-up
	-go run cmd/main.go

logs:
	docker logs -f url-shortener-app

logs-all:
	docker compose logs -f