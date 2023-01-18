package authserver

import (
	"context"
	"log"
	"net"
	"net/http"
	"testing"

	"ecommerce/auth-service/pkg/delivery/grpc/pb"
	"ecommerce/auth-service/pkg/domain"
	"ecommerce/auth-service/pkg/domain/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func Test_Register(t *testing.T) {
	mockUsecase := new(mocks.AuthUsecase)
	mockResponse := domain.RegisterResponse{
		Status: http.StatusCreated,
	}
	mockUsecase.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(mockResponse)

	ctx := context.Background()
	conn := initializeGrpcServer(ctx, mockUsecase)
	defer conn.Close()

	request := &pb.RegisterRequest{
		Email:    "xxx@zz.com",
		Password: "password",
	}

	client := pb.NewAuthServiceClient(conn)
	response, _ := client.Register(ctx, request)

	assert.Equal(t, mockResponse.Status, response.Status)
}

func Test_Login(t *testing.T) {
	mockUsecase := new(mocks.AuthUsecase)
	mockResponse := domain.LoginResponse{
		Status: http.StatusOK,
	}
	mockUsecase.On("Login", mock.Anything, mock.AnythingOfType("*domain.User")).Return(mockResponse)

	ctx := context.Background()
	conn := initializeGrpcServer(ctx, mockUsecase)
	defer conn.Close()

	request := &pb.LoginRequest{
		Email:    "xxx@zz.com",
		Password: "password",
	}

	client := pb.NewAuthServiceClient(conn)
	response, _ := client.Login(ctx, request)

	assert.Equal(t, mockResponse.Status, response.Status)
}

func Test_Validate(t *testing.T) {
	mockUsecase := new(mocks.AuthUsecase)
	mockResponse := domain.ValidateResponse{
		Status: http.StatusOK,
	}
	mockUsecase.On("Validate", mock.Anything, mock.AnythingOfType("*domain.User")).Return(mockResponse)

	ctx := context.Background()
	conn := initializeGrpcServer(ctx, mockUsecase)
	defer conn.Close()

	request := &pb.ValidateRequest{
		Token: "xxx",
	}

	client := pb.NewAuthServiceClient(conn)
	response, _ := client.Validate(ctx, request)

	assert.Equal(t, mockResponse.Status, response.Status)
}

func initializeGrpcServer(ctx context.Context, usecase *mocks.AuthUsecase) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer(usecase)))
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func dialer(usecase *mocks.AuthUsecase) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()

	authServer := NewGrpcAuthServer(usecase)

	pb.RegisterAuthServiceServer(grpcServer, authServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
