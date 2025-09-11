.PHONY: test

test-unit:
	go test -race $(shell go list ./... | grep -v /test | grep -v /pb) -coverprofile=coverage.out
cover:
	go tool cover -html=coverage.out
protos:
	cd proto/usermanagementapi/v1 && buf generate --clean
mocks:
	go install github.com/vektra/mockery/v2@latest
	mockery --dir=clients/ports --all --output=test/mocks