FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/ecommerce/auth-service
COPY . $GOPATH/src/ecommerce/auth-service
RUN go build -o auth-service ./cmd/main.go

EXPOSE 50051
ENTRYPOINT ["./auth-service"]