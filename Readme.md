[migrate tool]
export POSTGRESQL_URL='postgres://selectel:selectel@127.0.0.1:5432/selectel?sslmode=disable'

./migrate -path ./../migrations/ -database ${POSTGRESQL_URL} up