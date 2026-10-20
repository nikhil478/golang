how to run this appplication is by go run main.go

Docker command 

docker run --name go-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=testdb \
  -p 5432:5432 \
  -d postgres:alpine


Initial script
docker exec -it go-postgres psql -U postgres -d testdb -c "
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);"


Autocanon for testing

npm install -g autocannon

autocannon -c 100 -d 10 http://localhost:8080/items

autocannon -c 100 -d 10 \
  -m POST \
  -H "Content-Type: application/json" \
  -b '{"name":"Performance Testing"}' \
  http://localhost:8080/items
