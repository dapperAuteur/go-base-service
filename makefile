SHELL := /bin/bash

# =======================================
# Building containers

all:
	service-api

service-api:
	docker build \
		-f zarf/docker/dockerfile.service-api \
		-t service-api-amd64:1.0 \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within k8s/dev

kind-up:
	kind create cluster --image kindest/node:v1.20.2 --name awe-ful-starter-cluster --config zarf/k8s/dev/kind-config.yaml

kind-down:
	kind delete cluster --name awe-ful-starter-cluster
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