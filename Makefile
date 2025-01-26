all: postgres createdb

postgres:
	docker run --name some-postgres -p 5432:5432 -e POSTGRES_PASSWORD=qwerty -e POSTGRES_USER=root -d postgres:12-alpine

del_postgres:
	docker stop some-postgres
	docker rm some-postgres

createdb:
	@container_name="some-postgres"; \
	container_id=$$(docker ps -q -f name=$$container_name); \
	if [ -n "$$container_id" ]; then \
		attempt=1; \
		until docker exec some-postgres psql -U root -c '\q' 2>/dev/null; do \
			printf "Waiting for PostgreSQL to be ready (attempt #$$attempt)... \r"; \
			attempt=$$((attempt+1)); \
			sleep 1; \
		done; \
		printf "\n"; \
		echo "postgres container is ready!"; \
		docker exec -it some-postgres createdb --username=root --owner=root simplebank; \
	else \
		echo "Container $$container_name is not running"; \
	fi
	
dropdb:
	docker exec -it some-postgres dropdb simplebank

migrateup:
	migrate -path db/migration/ -database "postgresql://root:qwerty@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration/ -database "postgresql://root:qwerty@localhost:5432/simplebank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test ./... -v -cover

clean: dropdb del_postgres

.PHONY: createdb dropdb postgres del_postgres migrateup migratedown sqlc clean

