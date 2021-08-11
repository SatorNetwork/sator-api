# Image with necessary dependencies
FROM golang:alpine AS container
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh curl ca-certificates
RUN mkdir -p /src
WORKDIR /src
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download


# Go application builder
FROM container AS builder
WORKDIR /src
COPY . .
ARG LATEST_COMMIT=undefined_commit
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix nocgo -o /bin/migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix nocgo \
    -ldflags "-X main.buildTag=`date -u +%Y%m%d.%H%M%S`-${LATEST_COMMIT}" \
    -o /bin/api ./cmd/api


# Run go application
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
ARG APP_PORT=8080
RUN mkdir -p /migrations
COPY --from=builder /bin/* /
COPY --from=builder /src/migrations/* /migrations
COPY --from=builder /src/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE ${APP_PORT}
ENTRYPOINT [ "/entrypoint.sh" ]