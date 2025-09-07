.PHONY: test

test-unit:
	go test -race $(shell go list ./... | grep -v /test | grep -v /pb) -coverprofile=coverage.out
cover:
	go tool cover -html=coverage.out
test-integration:
	go test -race test/integration/*.go
swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/main.go -o app/docs