.SILENT:

run:
	go run cmd/web/main.go

run-worker:
	go run cmd/worker/main.go

test:
	go test -v ./test/

migrate:
	migrate -database "mysql://root:@tcp(localhost:3306)/golang_clean_architecture?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations up

migrate-new:
	echo "please run: migrate create -ext sql -dir db/migrations create_table_xxx"

generate:
	go generate ./...

swag:
	swag fmt && swag init --parseDependency --parseInternal --generalInfo ./cmd/web/main.go --output ./api/

docker-compose:
	docker compose up
