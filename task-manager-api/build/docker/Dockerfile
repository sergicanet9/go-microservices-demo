FROM golang:alpine AS builder

WORKDIR /opt/go-microservices-demo/task-manager-api

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o bin/main cmd/main.go

FROM alpine:latest

COPY --from=builder /opt/go-microservices-demo/task-manager-api/bin/main /opt/go-microservices-demo/task-manager-api/bin/main
COPY --from=builder /opt/go-microservices-demo/task-manager-api/config/*.json /opt/go-microservices-demo/task-manager-api/config/

WORKDIR /opt/go-microservices-demo/task-manager-api

ARG version
ENV v=$version

ARG environment
ENV env=$environment

ARG httpPort
ENV hp=$httpPort

ARG dsn
ENV dsn=$dsn

CMD ["sh", "-c", "bin/main --ver $v --env $env --hport $hp --dsn $dsn --jsecret $jsecret"]
