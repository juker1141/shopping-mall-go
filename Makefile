DB_URL=postgresql://postgres:password@localhost:5432/shopping_mall?sslmode=disable

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:15.3-alpine

createdb:
	docker exec -it postgres createdb -U postgres -O postgres shopping_mall

dropdb:
	docker exec -it postgres dropdb -U postgres shopping_mall

migrateup:
	migrate -path db/migration -database "${DB_URL}" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "${DB_URL}" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/juker1141/shopping-mall-go/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown new_migration sqlc test server mock