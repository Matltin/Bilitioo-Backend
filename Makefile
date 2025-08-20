postgres:
	docker run --name test_for_bilitioo_db --network bilitioo-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

psql:
	docker exec -it test_for_bilitioo_db psql -U root -d bilitioo

createdb:
	docker exec -it test_for_bilitioo_db createdb --username=root --owner=root bilitioo

dropdb:
	docker exec -it test_for_bilitioo_db dropdb bilitioo

migrateup:
	migrate -path db/migrate -database  "postgresql://root:secret@localhost:5432/bilitioo?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migrate -database  "postgresql://root:secret@localhost:5432/bilitioo?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrate -database "postgresql://root:secret@localhost:5432/bilitioo?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrate -database "postgresql://root:secret@localhost:5432/bilitioo?sslmode=disable" -verbose down 1

new_migrate:
	migrate create -ext sql -dir db/migrate -seq $(name)

dockerup:
	docker compose up -d

dockerdown:
	docker compose down

dockerlogs:
	docker compose logs -f

dockerstart:
	docker compose start

dockerstop:
	docker compose stop

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

restart:
	migrate -path db/migrate -database "postgresql://root:secret@localhost:5432/bilitioo?sslmode=disable" -verbose down
	migrate -path db/migrate -database  "postgresql://root:secret@localhost:5432/bilitioo?sslmode=disable" -verbose up

build:
	docker build -t bilitioo:latest . 

bilitioo:
	docker run --name bilitioo --network bilitioo-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@172.19.0.3:5432/bilitioo?sslmode=disable" bilitioo:latest

start:
	swag init
	swag fmt
	kill -9 $(shell lsof -t -i :3000) 2>/dev/null || true
	go run main.go

.PHONY: postgres psql createdb dropdb \
	migrateup migrateup1 migratedown migratedown1 \
	new_migrate dockerup dockerdown dockerlogs dockerstart dockerstop \
	sqlc test restart build bilitioo
