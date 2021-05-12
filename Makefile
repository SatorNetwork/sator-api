# Variables
LATEST_COMMIT := $$(git rev-parse HEAD)

.PHONY: help build docker

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
	-o ./app ./cmd/

docker: ## Build docker image
	@docker build -t satorapi:latest .
	@go clean

migrate: ## Run all migrations
	@rm -Rvf migrations/*.sql && \
	cp -Rvf ./**/repository/sql/migrations/*.sql migrations/ && \
	sql-migrate up -config=migrations/dbconfig.yml