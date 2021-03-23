SHELL := /bin/bash

run:
	go run app/kickball-api/main.go
tidy:
	go mod tidy
	go mod vendor