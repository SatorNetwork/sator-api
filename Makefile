# Variables
LATEST_COMMIT := $$(git rev-parse HEAD)

.PHONY: help build docker migrate up down

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
%:
	@:

build: ## Build the app
	@go clean
	@CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build \
	-a -installsuffix nocgo \
	-ldflags "-X main.buildTag=`date -u +%Y%m%d.%H%M%S`-$(LATEST_COMMIT)" \
	-o ./app ./cmd/api/main.go

run-local:
	@APP_PORT=8080 \
	APP_BASE_URL=https://aec45cb3e117.ngrok.io/ \
	DATABASE_URL=postgresql://pguser:pgpass@127.0.0.1/pgdb?sslmode=disable \
	JWT_SIGNING_KEY=secret \
	QUIZ_WS_CONN_URL=https://aec45cb3e117.ngrok.io/quiz \
	SOLANA_API_BASE_URL=https://api.devnet.solana.com/ \
	go run -ldflags "-X main.buildTag=`date -u +%Y%m%d.%H%M%S`-$(LATEST_COMMIT)" cmd/api/main.go

migrate: ## Run all migrations on server
	@rm -Rvf migrations/*.sql && \
	cp -Rvf ./svc/**/repository/sql/migrations/*.sql migrations/ && \
	./bin/migrate

migrate-local: ## Run all migrations on local environment
	@rm -Rvf migrations/*.sql && \
	cp -Rvf ./svc/**/repository/sql/migrations/*.sql migrations/ && \
	DATABASE_URL=postgresql://pguser:pgpass@127.0.0.1/pgdb?sslmode=disable \
	go run ./cmd/migrate/main.go

up: ## Run all needed containers, including postgres with exposed port :5432
	docker-compose up -d

down: ## Stop and remove all related containers
	docker-compose down -v --rmi=local


solana: ## Prepage solana accounts
	@DATABASE_URL=postgresql://pguser:pgpass@127.0.0.1/pgdb?sslmode=disable \
	go run cmd/solana/main.go