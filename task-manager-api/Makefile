.PHONY: test

test-unit:
	go test -race $(shell go list ./... | grep -v /test | grep -v /pb) -coverprofile=coverage.out
cover:
	go tool cover -html=coverage.out
swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/main.go -o app/docs
mocks:
	go install github.com/vektra/mockery/v2@latest
	mockery --dir=core/ports --all --output=test/mocks