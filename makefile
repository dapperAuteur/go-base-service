SHELL := /bin/bash

# =======================================
# Building containers

all: service-api

service-api: docker build \
			-f zarf/docker/dockerfile.service-api \
			-t service-api-amd64:1.0 \
			--build-arg VCS_REF=`git rev-parse HEAD` \
			--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
			.

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