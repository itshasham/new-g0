SWAG ?= swag

.PHONY: swagger swagger-install generate dev run test fmt

swagger: swagger-install
	$(SWAG) init -g main.go -o ./docs

swagger-install:
	@command -v $(SWAG) >/dev/null || go install github.com/swaggo/swag/cmd/swag@v1.16.4

generate:
	go generate ./...

dev: swagger
	go run main.go dev

run:
	go run main.go

test:
	go test ./...

fmt:
	gofmt -w .
