FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/ecommerce/order-service
COPY . $GOPATH/src/ecommerce/order-service
RUN go build -o order-service ./cmd/main.go

EXPOSE 50053
ENTRYPOINT ["./order-service"]