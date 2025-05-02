postgres:
	docker run --name postgres1 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it postgres1 createdb --username=root --owner=root simple_bank
	
dropdb:
	docker exec -it postgres1 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlcinit:
	docker run --rm -v "${shell cd}:/src" -w /src sqlc/sqlc init

sqlc:
	docker run --rm -v "${shell cd}:/src" -w /src sqlc/sqlc generate

test:
	go test -v -cover ./...

build:
	go build .
 
.PHONY: postgres createdb dropdb migrateup migratedown sqlcinit sqlcn test build migrateup1 migratedown1