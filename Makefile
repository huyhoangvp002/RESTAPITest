DB_URL=postgresql://root:secret@localhost:5432/Test_shop?sslmode=disable
postgres:
	sudo docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	sudo docker exec -it postgres12 createdb --username=root --owner=root Test_shop
dropdb:
	sudo docker exec -it postgres12 dropdb  Test_shop
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
sqlc:
	sqlc generate

mock:
	mockgen -package=mockdb github.com/huyhoangvp002/RESTAPITest/db/sqlc  Store  > db/mock/store.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc mock
