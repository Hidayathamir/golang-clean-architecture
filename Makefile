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

run-workeridempotencycleaner:
	mkdir -p logs
	$(RUN_CMD) cmd/workeridempotencycleaner/main.go >> logs/workeridempotencycleaner_log.jsonl 2>&1

go-test:
	$(TEST_CMD) -count=1 -v ./internal/... >> logs/go_test.jsonl 2>&1

go-e2e-test:
	$(TEST_CMD) -count=1 -v ./test/e2etest/... >> logs/go_e2e_test.jsonl 2>&1

#################################### 

migrate:
	$(RUN_CMD) cmd/migrate/main.go

#################################### 

docker-compose-up:
	docker compose up -d

docker-compose-up-clean:
	docker compose down -v && docker compose up -d

docker-compose-down:
	docker compose down

docker-compose-down-clean:
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

check-rules:
	go run ./scripts/checkrules

# Enforce Conventional Commits on a commit message file (see .githooks/commit-msg).
check-commit-msg:
	@if [ -z "$(MSG_FILE)" ]; then echo "usage: make check-commit-msg MSG_FILE=<file>"; exit 2; fi
	go run ./scripts/checkcommitmsg "$(MSG_FILE)"

# Render docs/*.mermaid to docs/*.png (requires Node; fetches mmdc via npx).
# -s 3 renders at 3x scale for crisp text on hi-dpi displays.
diagram:
	npx -y @mermaid-js/mermaid-cli -i docs/architecture.mermaid -o docs/architecture.png -b white -s 3 -p docs/puppeteer-config.json

# One-time setup after clone: activate the shared pre-commit hook.
setup:
	git config core.hooksPath .githooks
	@echo "✅ Git hooks activated (core.hooksPath = .githooks)"

# Fails if generated code is stale: regenerate, then assert no diff.
# Run in CI. Requires a clean working tree for the two generated dirs.
check-codegen:
	$(MAKE) generate
	$(MAKE) swag
	@git diff --exit-code -- internal/mock api || \
		(echo "❌ generated code is stale. Run 'make generate && make swag' and commit the result." && exit 1)

add-func-name:
	@find ./internal ./pkg ./cmd -name '*.go' -not -path '*/mock/*' -not -path '*/pkg/errkit/cmd/*' | xargs go run ./pkg/errkit/cmd/addfuncname

check-tools:
	@echo "🔍 Checking required tools..."
	@if command -v go >/dev/null 2>&1; then \
		echo "✔ go installed"; \
	else \
		echo "❌ go not found. Install: https://go.dev/"; \
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
	@if command -v npm >/dev/null 2>&1; then \
		echo "✔ npm installed"; \
	else \
		echo "❌ npm not found. Install Node.js: https://nodejs.org/"; \
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
