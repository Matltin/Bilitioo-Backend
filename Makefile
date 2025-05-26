postgres:
	docker run --name test_for_bilitioo_db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

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

.PHONY: sqlc test	