.PHONY: generate
MODULE := github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2
generate:
	protoc -I. \
	-I api/third_party \
	--go_out=. --go_opt=module=$(MODULE) \
	--validate_out=. --validate_opt=module=$(MODULE),lang=go \
	--go-grpc_out=. --go-grpc_opt=module=$(MODULE) \
	api/product/v1/product.proto

.PHONY: run-all
run-all:
	go run cmd/cart-service/main.go & \
	go run cmd/product-service/main.go & \
	wait

.PHONY: test
test:
	go test ./...

.PHONY: test-api
test-api:
	go test -v -count=1 -tags=api ./test/

.PHONY: test-race
test-race:
	go test -race ./...
