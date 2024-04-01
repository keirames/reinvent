docker run \
 --name test-db \
 -d \
 -e POSTGRES_USER=postgres \
 -e POSTGRES_PASSWORD=12345678 \
 -e POSTGRES_DB=recurrence \
 -p 5432:5432 \
 postgres

export POSTGRESQL_URL='postgres://postgres:12345678@localhost:5432/recurrence?sslmode=disable'

migrate create -ext sql -dir db/migrations -seq create_users_table
migrate -database ${POSTGRESQL_URL} -path db/migrations up
