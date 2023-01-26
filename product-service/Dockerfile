FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/ecommerce/product-service
COPY . $GOPATH/src/ecommerce/product-service
RUN go build -o product-service ./cmd/main.go

EXPOSE 50052
ENTRYPOINT ["./product-service"]