.SILENT:

run:
	mkdir -p logs
	go run cmd/web/main.go >> logs/web_log.jsonl 2>&1

run-worker:
	mkdir -p logs
	go run cmd/worker/main.go >> logs/worker_log.jsonl 2>&1

run-cron:
	mkdir -p logs
	go run cmd/cron/main.go >> logs/cron_log.jsonl 2>&1

go-test:
	go test -count=1 -v ./internal/... >> logs/go_test.jsonl 2>&1

go-e2e-test:
	go test -count=1 -v ./test/e2etest/... >> logs/go_e2e_test.jsonl 2>&1

#################################### 

migrate:
	go run cmd/migrate/main.go

migrate-new:
	echo "please run: migrate create -ext sql -dir db/migrations create_table_xxx"

#################################### 

docker-compose:
	docker compose down -v && docker compose up

docker-validate:
	docker ps --format "{{.Names}}\t{{.Status}}"

#################################### 

run-clean:
	make clean && make run

clean:
	make generate && make swag && make format && echo "done"

format:
	golangci-lint run ./... --fix --tests=false

generate:
	rm -rf internal/mock &&	go generate ./internal/...

swag:
	rm -rf api/ && swag fmt --exclude ./internal/mock && swag init --parseDependency --parseInternal --generalInfo ./cmd/web/main.go --output ./api/

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
	@if command -v moq >/dev/null 2>&1; then \
		echo "✔ moq installed"; \
	else \
		echo "❌ moq not found. Install: https://github.com/matryer/moq"; \
	fi
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "✔ golangci-lint installed"; \
	else \
		echo "❌ golangci-lint not found. Install: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"; \
	fi
	@echo "✅ Done checking tools."

OLD_MODULE := github.com/Hidayathamir/golang-clean-architecture

# Detect OS (Darwin = macOS, Linux = Linux/WSL)
OS := $(shell uname)

rename-go-mod:
	@read -p "👉 Enter new Go module name: " NEW_MODULE; \
	echo "🔄 Setting module to $$NEW_MODULE"; \
	if [ "$(OS)" = "Darwin" ]; then \
		grep -rl "$(OLD_MODULE)" . | xargs sed -i '' "s|$(OLD_MODULE)|$$NEW_MODULE|g"; \
	else \
		grep -rl "$(OLD_MODULE)" . | xargs sed -i "s|$(OLD_MODULE)|$$NEW_MODULE|g"; \
	fi; \
	echo "⚙️ Running 'make generate'"; \
	$(MAKE) generate; \
	echo "⚙️ Running 'make swag'"; \
	$(MAKE) swag
