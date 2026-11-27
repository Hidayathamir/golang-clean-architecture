.SILENT:

# tuistory wraps the long-running servers in named sessions so rerunning `make run-*`
# restarts them idempotently instead of spawning a duplicate (port/consumer conflicts).
# It is REQUIRED for the run-* targets. One-shot targets (migrate, tests) use plain
# go run/go test so they exit cleanly without a TUI.
HAS_TUISTORY := $(shell command -v tuistory 2>/dev/null)

run-webserver:
	@$(if $(HAS_TUISTORY),, $(error ❌ tuistory not found. Install: npm install -g tuistory (see 'make check-tools'). The run-* targets need it for idempotent restarts.))
	mkdir -p logs
	tuistory -- go run cmd/webserver/main.go >> logs/webserver_log.jsonl 2>&1

run-workerconsumer:
	@$(if $(HAS_TUISTORY),, $(error ❌ tuistory not found. Install: npm install -g tuistory (see 'make check-tools'). The run-* targets need it for idempotent restarts.))
	mkdir -p logs
	tuistory -- go run cmd/workerconsumer/main.go >> logs/workerconsumer_log.jsonl 2>&1

run-workerproducer:
	@$(if $(HAS_TUISTORY),, $(error ❌ tuistory not found. Install: npm install -g tuistory (see 'make check-tools'). The run-* targets need it for idempotent restarts.))
	mkdir -p logs
	tuistory -- go run cmd/workerproducer/main.go >> logs/workerproducer_log.jsonl 2>&1

run-workeridempotencycleaner:
	@$(if $(HAS_TUISTORY),, $(error ❌ tuistory not found. Install: npm install -g tuistory (see 'make check-tools'). The run-* targets need it for idempotent restarts.))
	mkdir -p logs
	tuistory -- go run cmd/workeridempotencycleaner/main.go >> logs/workeridempotencycleaner_log.jsonl 2>&1

go-test:
	mkdir -p logs
	go test -count=1 -v ./internal/... >> logs/go_test.jsonl 2>&1

go-e2e-test:
	mkdir -p logs
	go test -count=1 -v ./test/e2etest/... >> logs/go_e2e_test.jsonl 2>&1

#################################### 

migrate:
	go run cmd/migrate/main.go

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

# One-time onboarding after clone. Idempotent — safe to rerun.
# 1) verify tools (hard-fails if missing)  2) download Go deps
# 3) activate git hooks  4) verify the hook is live  5) print next steps.
setup:
	@echo "🚀 Setting up the project..."
	@$(MAKE) --no-print-directory check-tools
	@echo "📦 Downloading Go dependencies..."
	@go mod download
	@echo "🪝 Activating git hooks..."
	@git config core.hooksPath .githooks
	@test "$$(git config core.hooksPath)" = ".githooks" || (echo "❌ Failed to activate git hooks." && exit 1)
	@if [ ! -x .githooks/pre-commit ]; then echo "❌ .githooks/pre-commit is not executable."; exit 1; fi
	@echo "✅ Git hooks activated (core.hooksPath = .githooks)"
	@echo ""
	@echo "✅ Setup complete. Next steps:"
	@echo "   1. make docker-compose-up            # start infra (wait for it to finish)"
	@echo "   2. make migrate                      # apply DB migrations"
	@echo "   3. make run-webserver                # and the other make run-* servers"

# Fails if generated code is stale: regenerate, then assert no diff.
# Runs in the pre-push hook (.githooks/pre-push); requires a clean working tree
# for the two generated dirs.
check-codegen:
	$(MAKE) generate
	$(MAKE) swag
	@git diff --exit-code -- internal/mock api || \
		(echo "❌ generated code is stale. Run 'make generate && make swag' and commit the result." && exit 1)

add-func-name:
	@find ./internal ./pkg ./cmd -name '*.go' -not -path '*/mock/*' -not -path '*/pkg/errkit/cmd/*' | xargs go run ./pkg/errkit/cmd/addfuncname

check-tools:
	@echo "🔍 Checking required tools..."
	@missing=0; \
	if command -v go >/dev/null 2>&1; then echo "✔ go installed"; else echo "❌ go not found. Install: https://go.dev/"; missing=1; fi; \
	if command -v docker >/dev/null 2>&1; then echo "✔ docker installed"; else echo "❌ docker not found. Install: https://www.docker.com/"; missing=1; fi; \
	if command -v swag >/dev/null 2>&1; then echo "✔ swag installed"; else echo "❌ swag not found. Install: go install github.com/swaggo/swag/cmd/swag@latest"; missing=1; fi; \
	if command -v moq >/dev/null 2>&1; then echo "✔ moq installed"; else echo "❌ moq not found. Install: go install github.com/matryer/moq@latest"; missing=1; fi; \
	if command -v golangci-lint >/dev/null 2>&1; then echo "✔ golangci-lint installed"; else echo "❌ golangci-lint not found. Install: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"; missing=1; fi; \
	if command -v npm >/dev/null 2>&1; then echo "✔ npm installed"; else echo "❌ npm not found. Install Node.js: https://nodejs.org/"; missing=1; fi; \
	if command -v tuistory >/dev/null 2>&1; then echo "✔ tuistory installed"; else echo "❌ tuistory not found. Required for 'make run-*'. Install: npm install -g tuistory"; missing=1; fi; \
	if [ $$missing -ne 0 ]; then echo "❌ One or more required tools are missing — install them, then rerun 'make setup'."; exit 1; fi; \
	echo "✅ All required tools installed."

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
