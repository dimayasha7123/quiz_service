POSTGRES_HOST=localhost
POSTGRES_PORT=5433
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=quiz_tg_bot

GOOSE_DRIVER=postgres
#GOOSE_DBSTRING="host=${POSTGRES_HOST} port=${POSTGRES_PORT} user=${POSTGRES_USER} password=${POSTGRES_PASSWORDP} dbname=${POSTGRES_DB} sslmode=disable"
GOOSE_DBSTRING="${GOOSE_DRIVER}://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}"

#GOOSE_DRIVER?=postgres
#GOOSE_DBSTRING?="host=localhost port=5433 user=postgres password=postgres dbname=quiz_tg_bot sslmode=disable"
#

goose_status:
	goose -dir ./migrations ${GOOSE_DRIVER} ${GOOSE_DBSTRING} status

goose_up:
	goose -dir ./migrations ${GOOSE_DRIVER} ${GOOSE_DBSTRING} up

goose_reset:
	goose -dir ./migrations ${GOOSE_DRIVER} ${GOOSE_DBSTRING} reset

init_postgres:
	sudo docker run --name tg_bot_postgres -p 5433:5432 \
	-e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e \
	POSTGRES_DB=quiz_tg_bot -d postgres

annihilate_postgres:
	sudo docker container stop tg_bot_postgres
	sudo docker container rm tg_bot_postgres

start:
	go run ./cmd/client/main.go