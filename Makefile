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

# Rem protoc --go_out=pkg --go_opt=paths=source_relative --go-grpc_out=pkg --go-grpc_opt=paths=source_relative api/api.proto
# Rem protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=. --grpc-gateway_opt generate_unbound_methods=true --openapiv2_out . api.proto
# protoc -I ./api --go_out ./pkg/api --go_opt paths=source_relative --go-grpc_out ./pkg/api --go-grpc_opt paths=source_relative --grpc-gateway_out ./pkg/api --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --openapiv2_out ./pkg/api api/api.proto