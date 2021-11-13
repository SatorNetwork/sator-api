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

migrate: ## Run all migrations on server
	@rm -Rvf migrations/*.sql && \
	cp -Rvf ./svc/**/repository/sql/migrations/*.sql migrations/ && \
	./bin/migrate

up: ## Run all needed containers, including postgres with exposed port :5432
	docker-compose up -d

down: ## Stop and remove all related containers
	docker-compose down -v --rmi=local

test:
	go test ./internal/test/...
