package orderserver

import (
	"context"
	"ecommerce/order-service/pkg/delivery/grpc/pb"
	"ecommerce/order-service/pkg/domain"
	"ecommerce/order-service/pkg/domain/mocks"
	"errors"
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

const (
	no_product_and_has_err_product_id_is_1 = 1
	no_product_without_err_product_id_is_2 = 2
	size_is_not_confitable_product_id_is_3 = 3
	conflict_and_has_error_product_id_is_4 = 4
	conflict_without_error_product_id_is_5 = 5
	all_ok_o_product_order_product_id_is_6 = 6
)

type MockGrpcProductService struct {
	pb.UnimplementedProductServiceServer
}

func (m *MockGrpcProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	return &pb.CreateProductResponse{}, errors.New("test")
}

func (m *MockGrpcProductService) FindOne(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {
	data := &pb.FindOneData{
		Id:    req.Id,
		Name:  "",
		Price: 2,
	}

	switch req.Id {
	case no_product_and_has_err_product_id_is_1:
		return &pb.FindOneResponse{Status: http.StatusNotFound}, errors.New("test")
	case no_product_without_err_product_id_is_2:
		return &pb.FindOneResponse{Status: http.StatusNotFound}, nil

	case size_is_not_confitable_product_id_is_3:
		data.Stock = 1
		return &pb.FindOneResponse{Status: http.StatusOK, Error: "", Data: data}, nil

	case conflict_and_has_error_product_id_is_4,
		conflict_without_error_product_id_is_5,
		all_ok_o_product_order_product_id_is_6:

		data.Stock = 5
		return &pb.FindOneResponse{Status: http.StatusOK, Error: "", Data: data}, nil

	}

	return &pb.FindOneResponse{Status: http.StatusNotFound}, errors.New("test")
}

func (m *MockGrpcProductService) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
	switch req.Id {
	case conflict_and_has_error_product_id_is_4:
		return &pb.DecreaseStockResponse{Status: http.StatusConflict, Error: "failed DecreaseStock"}, nil

	case conflict_without_error_product_id_is_5:
		return &pb.DecreaseStockResponse{Status: http.StatusConflict, Error: ""}, nil

	case all_ok_o_product_order_product_id_is_6:
		return &pb.DecreaseStockResponse{Status: http.StatusOK, Error: ""}, nil
	}

	return &pb.DecreaseStockResponse{Status: http.StatusBadRequest, Error: "not empy"}, nil
}

func Test_CreateProduct(t *testing.T) {
	mockUsecase := new(mocks.OrderUsecase)
	productclient := NewMockProductCient()
	mockResponse := domain.OrderResponse{
		Status: http.StatusOK,
		Error:  "",
	}
	mockUsecase.On("Create", mock.Anything, mock.AnythingOfType("*domain.Order")).
		Return(&mockResponse)
	mockUsecase.On("Delete", mock.Anything, mock.AnythingOfType("int64")).
		Return(&mockResponse)

	ctx := context.Background()
	cc, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(productclient, mockUsecase)))
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	client := pb.NewOrderServiceClient(cc)

	requests := []struct {
		Name           string
		ExpectedStatus int64
		request        *pb.CreateOrderRequest
	}{
		{Name: "no_product_and_has_err_product_id_is_1", ExpectedStatus: http.StatusBadRequest, request: &pb.CreateOrderRequest{ProductId: no_product_and_has_err_product_id_is_1, Quantity: 1, UserId: 1}},
		{Name: "no_product_without_err_product_id_is_2", ExpectedStatus: http.StatusNotFound, request: &pb.CreateOrderRequest{ProductId: no_product_without_err_product_id_is_2, Quantity: 1, UserId: 2}},
		{Name: "size_is_not_confitable_product_id_is_3", ExpectedStatus: http.StatusConflict, request: &pb.CreateOrderRequest{ProductId: size_is_not_confitable_product_id_is_3, Quantity: 2, UserId: 3}},
		{Name: "conflict_and_has_error_product_id_is_4", ExpectedStatus: http.StatusConflict, request: &pb.CreateOrderRequest{ProductId: conflict_and_has_error_product_id_is_4, Quantity: 3, UserId: 4}},
		{Name: "conflict_without_error_product_id_is_5", ExpectedStatus: http.StatusConflict, request: &pb.CreateOrderRequest{ProductId: conflict_without_error_product_id_is_5, Quantity: 3, UserId: 5}},
		{Name: "all_ok_o_product_order_product_id_is_6", ExpectedStatus: http.StatusCreated, request: &pb.CreateOrderRequest{ProductId: all_ok_o_product_order_product_id_is_6, Quantity: 3, UserId: 6}},
	}

	for _, tt := range requests {
		t.Run(tt.Name, func(t *testing.T) {
			response, _ := client.CreateOrder(ctx, tt.request)
			assert.Equal(t, int64(tt.ExpectedStatus), response.Status)
		})
	}
}

func NewMockProductCient() *grpcProductClient {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialerProductServer()))
	if err != nil {
		log.Fatal(err)
	}

	return NewGrpcProductClient(conn)
}

func dialerProductServer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, &MockGrpcProductService{})

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func dialer(pclient *grpcProductClient, ucase domain.OrderUsecase) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	orderServer := NewGrpcOrderServer(ucase, pclient)

	pb.RegisterOrderServiceServer(grpcServer, orderServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
