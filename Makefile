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

docker: ## Build docker image
	rm -Rvf migrations/*.sql && cp -Rvf ./svc/**/repository/sql/migrations/*.sql migrations/ \
	&& docker build -f Dockerfile --build-arg LATEST_COMMIT=$(LATEST_COMMIT) -t satorapi:latest . \
	&& docker scan satorapi:latest

run-local: ## Run api via `go run`
	@APP_PORT=8080 \
	APP_BASE_URL=https://aec45cb3e117.ngrok.io/ \
	DATABASE_URL=postgresql://pguser:pgpass@127.0.0.1/pgdb?sslmode=disable \
	JWT_SIGNING_KEY=secret \
	QUIZ_WS_CONN_URL=https://aec45cb3e117.ngrok.io/quiz \
	SOLANA_API_BASE_URL=https://api.devnet.solana.com/ \
	POSTMARK_SERVER_TOKEN=local \
	POSTMARK_ACCOUNT_TOKEN=local \
	STORAGE_KEY=XXXXXXXXXX \
    STORAGE_SECRET=XXXXXXXXXX \
    STORAGE_ENDPOINT=https://nyc3.digitaloceanspaces.com \
    STORAGE_REGION=nyc3 \
    STORAGE_BUCKET=sator-media-storage \
    STORAGE_URL=https://sator-media-storage.nyc3.digitaloceanspaces.com \
    STORAGE_FORCE_PATH_STYLE=false \
    STORAGE_DISABLE_SSL=true \
    FIREBASE_BASE_URL=https://satorio.page.link \
    FIREBASE_WEB_API_KEY=XXXXXXXXXXXXXXXXXXXXXX \
    FIREBASE_MAIN_SITE_LINK=https://sator.io/ \
    FIREBASE_ANDROID_PACKAGE_NAME=com.satorio.app \
    FIREBASE_IOS_BUNDLE_ID=io.sator \
    FIREBASE_SUFFIX_OPTION=UNGUESSABLE \
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

request-airdrop-testnet:
	SOLANA_API_BASE_URL=https://api.testnet.solana.com \
	go run ./cmd/solana/airdrop/main.go

create-missed-wallets-testnet:
	DATABASE_URL=postgresql://sator-db-prod:y4ttotjayg3h141c@app-6c60d8cd-32e4-432b-86d6-c19c0a2c3052-do-user-9696296-0.b.db.ondigitalocean.com:25060/sator-db-prod?sslmode=require \
	SOLANA_API_BASE_URL=https://api.testnet.solana.com \
	go run ./cmd/solana/missedwallets/main.go

create-missed-wallets-devnet:
	DATABASE_URL=postgresql://db:zk131mmy8sdf93f5@app-03ef8dbc-d361-4581-a833-2671c70adb58-do-user-9696296-0.b.db.ondigitalocean.com:25060/db?sslmode=require \
	SOLANA_API_BASE_URL=https://api.devnet.solana.com \
	go run ./cmd/solana/missedwallets/main.go