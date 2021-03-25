SHELL := /bin/bash

run:
	go run app/service-api/main.go

runadmin:
	go run app/service-admin/main.go

test:
	go test -v ./... -count=1
	staticcheck ./...
	
tidy:
	go mod tidy
	go mod vendor