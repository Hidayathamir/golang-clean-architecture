.SILENT:

HAS_TUISTORY := $(shell command -v tuistory 2>/dev/null)
ifeq ($(HAS_TUISTORY),)
  RUN_CMD = go run
  TEST_CMD = go test
else
  RUN_CMD = tuistory -- go run
  TEST_CMD = tuistory -- go test
endif

run-webserver:
	mkdir -p logs
	$(RUN_CMD) cmd/webserver/main.go >> logs/webserver_log.jsonl 2>&1

run-workerconsumer:
	mkdir -p logs
	$(RUN_CMD) cmd/workerconsumer/main.go >> logs/workerconsumer_log.jsonl 2>&1

run-workerproducer:
	mkdir -p logs
	$(RUN_CMD) cmd/workerproducer/main.go >> logs/workerproducer_log.jsonl 2>&1

go-test:
	$(TEST_CMD) -count=1 -v ./internal/... >> logs/go_test.jsonl 2>&1

go-e2e-test:
	$(TEST_CMD) -count=1 -v ./test/e2etest/... >> logs/go_e2e_test.jsonl 2>&1

#################################### 

migrate:
	$(RUN_CMD) cmd/migrate/main.go

new-migration:
	sql-migrate new -config=dbconfig.yml -env=local

migrate-status:
	sql-migrate status -config=dbconfig.yml -env=local

migrate-up:
	sql-migrate up -config=dbconfig.yml -env=local

migrate-down:
	sql-migrate down -limit=0 -config=dbconfig.yml -env=local

#################################### 

docker-compose-up:
	docker compose down -v && docker compose up -d

docker-compose-down:
	docker compose down -v

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
	rm -rf api/ && swag fmt --exclude ./internal/mock && swag init --parseDependency --parseInternal --generalInfo ./cmd/webserver/main.go --output ./api/

add-func-name:
	@find ./internal ./pkg ./cmd -name '*.go' -not -path '*/mock/*' -not -path '*/pkg/errkit/cmd/*' | xargs go run ./pkg/errkit/cmd/addfuncname

check-tools:
	@echo "🔍 Checking required tools..."
	@if command -v go >/dev/null 2>&1; then \
		echo "✔ go installed"; \
	else \
		echo "❌ go not found. Install: https://go.dev/"; \
	fi
	@if command -v sql-migrate >/dev/null 2>&1; then \
		echo "✔ sql-migrate installed"; \
	else \
		echo "❌ sql-migrate not found. Install: go install github.com/rubenv/sql-migrate/sql-migrate@latest"; \
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
