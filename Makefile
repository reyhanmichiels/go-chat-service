.PHONY: swag-install
swag-install:
	@go install github.com/swaggo/swag/cmd/swag@v1.6.7

.PHONY: swaggo
swaggo:
	@/bin/rm -rf ./docs/swagger
	@`go env GOPATH`/bin/swag init -g ./src/cmd/main.go -o ./docs/swagger --parseInternal

.PHONY: build
build:
	@go build -o ./build/app ./src/cmd

.PHONY: build-alpine
build-alpine:
	@go mod tidy && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/app ./src/cmd

.PHONY: run
run: swaggo build
	@./build/app

.PHONY: run-tests
run-tests:
	@go clean -cache
	@go test -v -failfast `go list ./... | grep -i 'business'` -cover

.PHONY: lint
lint:
	@`go env GOPATH`/bin/golangci-lint run

.PHONY: lint-fast
lint-fast:
	@`go env GOPATH`/bin/golangci-lint run --fast

.PHONY: lint-install
lint-install:
	@./utils/lint_install

.PHONY: mock-install
mock-install:
	@go install go.uber.org/mock/mockgen@v0.4.0

.PHONY: mock
mock:
	@`go env GOPATH`/bin/mockgen -source src/business/domain/$(domain)/$(domain).go -destination src/business/domain/mock/$(domain)/$(domain).go