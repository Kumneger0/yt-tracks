project_name?=yt-tracks
BACKEND_STAMP := build/.backend.stamp
BACKEND_SOURCES := $(shell find grpc_server -path 'grpc_server/gen' -prune -o -name '__pycache__' -prune -o -type f -print 2>/dev/null) proto/music.proto pyproject.toml $(wildcard uv.lock) buf.gen.python.yaml scripts/build-backends.sh
GO_SOURCES := $(shell find . -path './.venv' -prune -o -path './build' -prune -o -path './docs' -prune -o -path './main.build' -prune -o -path './main.dist' -prune -o -name '*.go' -print 2>/dev/null) go.mod go.sum
ifeq ($(and $(wildcard internal/backend/binaries/wheel/*.whl),$(wildcard internal/backend/binaries/linux_amd64/uv),$(wildcard internal/backend/binaries/darwin_arm64/uv),$(wildcard internal/backend/binaries/windows_amd64/uv.exe)),)
.PHONY: $(BACKEND_STAMP)
endif

default: help

.PHONY: help
help: ## show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: build
build: $(project_name) ## build the Go application

$(project_name): $(BACKEND_STAMP) $(GO_SOURCES)
	@echo "--> Building Go application..."
	@go build -ldflags "-X main.version=$(shell git describe --abbrev=0 --tags) -X main.Debug=true" -o $(project_name)


.PHONY: install
install: build ## build and install yt-tracks to /usr/local/bin
	@echo "--> Installing yt-tracks to /usr/local/bin..."
	@sudo cp $(project_name) /usr/local/bin/
	@echo "--> Installation complete. Run 'yt-tracks' to start."

.PHONY: run
run: build ## build and run yt-tracks
	@./$(project_name)

.PHONY: bootstrap	
bootstrap: ## bootstrap go tools
	go generate -tags tools tools/tools.go

.PHONY: test
test: clean ## run go tests
	go test --cover -parallel=1 -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | sort -rnk3

.PHONY: cover
cover: ## display go test coverage
	go test -v -race $(shell go list ./... | grep -v /vendor/) -v -coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: fmt
fmt: ## format go files
	gofumpt -w .
	gci write .

.PHONY: lint
lint: ## lint go files
	golangci-lint run -c .golangci.yml

.PHONY: pre-commit
pre-commit:	## run pre-commit hooks
	pre-commit run --all-files

.PHONY: hooks
hooks: ## install git commit-msg hook for commitlint (local)
	@chmod +x scripts/hooks/commit-msg
	@git config core.hooksPath scripts/hooks
	@echo "--> Git hooks installed (commit-msg)."

.PHONY: proto
proto:
	@buf generate --template buf.gen.go.yaml
	@buf generate --template buf.gen.python.yaml
	@echo "compiled proto files successfully."

.PHONY: server-watch
server-watch: 
	@echo "Watching proto/music.proto for changes..."
	nodemon --watch proto/music.proto --exec "make proto"

.PHONY: server-run
server-run: ## run the python Connect RPC server
	@echo "Starting Connect RPC server..."
	@nodemon --watch grpc_server -e py --exec ".venv/bin/python -m grpc_server.main"

.PHONY: server-login
server-login: ## spin up http server for login
	@echo "Starting login flow..."
	.venv/bin/python grpc_server/main.py --login

.PHONY: server-sync
server-sync: ## sync python virtual environment dependencies
	@echo "Syncing virtual environment dependencies using uv..."
	uv sync

$(BACKEND_STAMP): $(BACKEND_SOURCES)
	@echo "--> Building Python backend assets"
	@./scripts/build-backends.sh
	@mkdir -p $(dir $@)
	@touch $@

.PHONY: backend
backend: ## build the Python backend wheel and fetch uv binaries for embedding
	@echo "--> Building Python backend assets"
	@./scripts/build-backends.sh
	@mkdir -p $(dir $(BACKEND_STAMP))
	@touch $(BACKEND_STAMP)

.PHONY: clean
clean: ## clean up both go and python generated files
	@echo "--> Cleaning up..."
	@rm -rf coverage.out dist/ $(project_name) $(BACKEND_STAMP)
	@rm -rf gen/*.go gen/genconnect
	@rm -rf grpc_server/gen/