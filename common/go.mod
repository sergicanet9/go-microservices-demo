module github.com/sergicanet9/go-microservices-demo/common

go 1.24.3

replace github.com/sergicanet9/go-microservices-demo/task-manager-api => ../task-manager-api

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2
	github.com/sergicanet9/go-microservices-demo/task-manager-api v0.0.0-20250911095929-947944135f1f
	github.com/stretchr/testify v1.11.1
	google.golang.org/genproto/googleapis/api v0.0.0-20250908214217-97024824d090
	google.golang.org/grpc v1.75.1
	google.golang.org/protobuf v1.36.9
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250826171959-ef028d996bc1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
