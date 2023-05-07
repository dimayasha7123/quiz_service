GOOSE_DRIVER?=postgres
GOOSE_DBSTRING?="host=localhost port=5432 user=postgres password=postgres dbname=quiz_service_db sslmode=disable"

goose_status:
	goose -dir ./migrations ${GOOSE_DRIVER} ${GOOSE_DBSTRING} status

goose_up:
	goose -dir ./migrations ${GOOSE_DRIVER} ${GOOSE_DBSTRING} up

goose_reset:
	goose -dir ./migrations ${GOOSE_DRIVER} ${GOOSE_DBSTRING} reset

init_postgres:
	sudo docker run --name testPostgres -p 5432:5432 \
	-e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e \
	POSTGRES_DB=quiz_service_db -d postgres

annihilate_postgres:
	sudo docker container stop testPostgres
	sudo docker container rm testPostgres

start:
	go run ./cmd/server/main.go