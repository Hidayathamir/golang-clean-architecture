.SILENT:

run:
	go run cmd/web/main.go

run-worker:
	go run cmd/worker/main.go

go-test:
	go test -count=1 -v ./internal/...

go-integration-test:
	go test -count=1 -v ./integrationtest/...

test:
	go test -count=1 ./...

#################################### 

migrate:
	go run cmd/migrate/main.go

migrate-new:
	echo "please run: migrate create -ext sql -dir db/migrations create_table_xxx"

#################################### 

docker-compose:
	docker compose down && docker compose up

docker-check:
	docker ps --format "{{.Names}}\t{{.Status}}"

#################################### 

run-clean:
	make clean && make run

clean:
	make generate && make swag && make format && echo "done"

format:
	golangci-lint run ./... --fix

generate:
	rm -rf internal/mock &&	go generate ./internal/...

swag:
	swag fmt --exclude ./internal/mock && swag init --parseDependency --parseInternal --generalInfo ./cmd/web/main.go --output ./api/

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
