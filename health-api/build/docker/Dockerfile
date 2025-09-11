FROM golang:alpine AS builder

WORKDIR /opt/go-microservices-demo/health-api

COPY ./common /opt/go-microservices-demo/common

COPY ./health-api/go.mod ./health-api/go.sum ./
RUN go mod download
COPY ./health-api .

RUN go build -o bin/main cmd/main.go

FROM alpine:latest

COPY --from=builder /opt/go-microservices-demo/health-api/bin/main /opt/go-microservices-demo/health-api/bin/main
COPY --from=builder /opt/go-microservices-demo/health-api/config/*.json /opt/go-microservices-demo/health-api/config/

WORKDIR /opt/go-microservices-demo/health-api

ARG version
ENV v=$version

ARG environment
ENV env=$environment

ARG httpPort
ENV hp=$httpPort

ARG tasksURL
ENV tasksurl=$tasksURL

ARG usersURL
ENV usersurl=$usersURL

CMD ["sh", "-c", "bin/main --ver $v --env $env --hport $hp --tasksurl $tasksurl --usersurl $usersurl"]