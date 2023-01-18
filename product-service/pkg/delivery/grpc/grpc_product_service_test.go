package productserver

import (
	"context"
	"ecommerce/product-service/pkg/delivery/grpc/pb"
	"ecommerce/product-service/pkg/domain"
	"ecommerce/product-service/pkg/domain/mocks"
	"log"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var mockUsecase *mocks.ProductUsecase

func Test_CreateProduct(t *testing.T) {
	mockUsecase = new(mocks.ProductUsecase)
	mockResponse := domain.ProdcutResponse{
		Status: http.StatusOK,
		Error:  "",
	}
	mockUsecase.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(mockResponse)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewProductServiceClient(conn)

	request := &pb.CreateProductRequest{
		Name:  "pc",
		Stock: 5,
		Price: 15000,
	}
	response, err := client.CreateProduct(ctx, request)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, int64(http.StatusOK), response.Status)
}

func Test_FindOne(t *testing.T) {
	mockUsecase = new(mocks.ProductUsecase)
	mockResponse := domain.ProdcutResponse{
		Status: http.StatusOK,
		Error:  "",
	}
	mockUsecase.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Product{}, mockResponse)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewProductServiceClient(conn)

	request := &pb.FindOneRequest{
		Id: 1,
	}
	response, err := client.FindOne(ctx, request)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, int64(http.StatusOK), response.Status)
}

func Test_DecreaseStock(t *testing.T) {
	mockUsecase = new(mocks.ProductUsecase)
	mockResponse := domain.ProdcutResponse{
		Status: http.StatusOK,
		Error:  "",
	}
	mockUsecase.On("DecreaseStock", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
		Return(mockResponse)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewProductServiceClient(conn)

	request := &pb.DecreaseStockRequest{
		Id:      1,
		OrderId: 2,
	}
	response, err := client.DecreaseStock(ctx, request)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, int64(http.StatusOK), response.Status)
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	productServer := NewGrpcProductService(mockUsecase)

	pb.RegisterProductServiceServer(grpcServer, productServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
