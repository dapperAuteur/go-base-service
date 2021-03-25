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

# kind-up-m1:
# 	kind create cluster --image rossgeorgiev/kind-node-arm64 --name awe-ful-starter-cluster --config zarf/k8s/dev/kind-config.yaml

kind-down:
	kind delete cluster --name awe-ful-starter-cluster

kind-load:
	kind load docker-image service-api-amd64:1.0 --name awe-ful-starter-cluster

kind-services:
	kustomize build zarf/k8s/dev | kubectl apply -f -

# kind-update: service
# 	kind load docker-image service-api-amd64:1.0 --name awe-ful-starter-cluster
# 	kubectl delete pods -lapp=service-api

kind-status:
	kubectl get nodes
	kubectl get pods --watch

kind-status-full:
	kubectl describe pod -lapp=service-api

kind-logs:
	kubectl logs -lapp=service-api --all-containers=true -f

# ==============================================
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