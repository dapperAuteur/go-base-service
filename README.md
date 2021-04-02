# go-base-service

# ==============================================================================
# Testing running system

# For testing a simple query on the system. Don't forget to `make seed` first.
# curl --user "admin@example.com:gophers" http://localhost:3000/v1/users/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1
# export TOKEN="COPY TOKEN STRING FROM LAST CALL"
# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users/1/2

# For testing load on the service.
# hey -m GET -c 100 -n 10000 -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users/1/2
# zipkin: http://localhost:9411
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,mem:memstats.Alloc"

# go install github.com/divan/expvarmon@latest

# // To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem

# ==============================================================================
# Building containers


==========
## Testing
The testing db uses `dbImage = "postgres:13-alpine"`.
It will need to be updated if the app db is changed.

Tests aren't working. Getting a reflect error:
```
--- FAIL: TestUser (25.00s)
panic: reflect: call of github.com/jmoiron/sqlx/reflectx.(*Mapper).TraversalsByNameFunc on string Value [recovered]
	panic: reflect: call of github.com/jmoiron/sqlx/reflectx.(*Mapper).TraversalsByNameFunc on string Value

goroutine 19 [running]:
```

run these commands after fixing test issue to see test coverage:
`go test -cover`
`go tool cover -html cover.out`

========
## Steps to get Docker, Kubernetes, and etc up and running
start docker
run the following commands in this order after the previous command has completed

* make kind-up
* make kind-load
* make kind-services
* make kind-status (should show DB up and running)