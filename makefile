SHELL := /bin/bash

run:
	go run app/kickball-api/main.go

runadmin:
	go run app/kickball-admin/main.go
tidy:
	go mod tidy
	go mod vendor