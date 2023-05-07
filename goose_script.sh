conn_string="host=${POSTGRES_HOST} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable"
echo "$conn_string"
goose postgres "${conn_string}" up
goose postgres "${conn_string}" status
exit 0