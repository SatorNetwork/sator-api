#!/bin/bash
go mod download
go build -o ./bin/api ./cmd/api/main.go
