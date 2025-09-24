.SILENT:

run:
	go run cmd/web/main.go

run-worker:
	go run cmd/worker/main.go

go-test:
	go test -v ./test/

migrate:
	go run cmd/migrate/main.go

migrate-new:
	echo "please run: migrate create -ext sql -dir db/migrations create_table_xxx"

generate:
	go generate ./internal/...

swag:
	swag fmt && swag init --parseDependency --parseInternal --generalInfo ./cmd/web/main.go --output ./api/

docker-compose:
	docker compose down && docker compose up

check-tools:
	@echo "🔍 Checking required tools..."
	@if command -v go >/dev/null 2>&1; then \
		echo "✔ go installed"; \
	else \
		echo "❌ go not found. Install: https://go.dev/"; \
	fi
	@if command -v migrate >/dev/null 2>&1; then \
		echo "✔ migrate installed"; \
	else \
		echo "❌ migrate not found. Install: https://github.com/golang-migrate/migrate"; \
	fi
	@if command -v docker >/dev/null 2>&1; then \
		echo "✔ docker installed"; \
	else \
		echo "❌ docker not found. Install: https://www.docker.com/"; \
	fi
	@if command -v swag >/dev/null 2>&1; then \
		echo "✔ swag installed"; \
	else \
		echo "❌ swag not found. Install: https://github.com/swaggo/swag"; \
	fi
	@echo "✅ Done checking tools."
