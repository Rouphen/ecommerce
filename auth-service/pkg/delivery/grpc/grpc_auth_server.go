package authserver

import (
	"context"
	"ecommerce/auth-service/pkg/delivery/grpc/pb"
	"ecommerce/auth-service/pkg/domain"
)

type grpcAuthServer struct {
	authUsecae domain.AuthUsecase
	pb.UnimplementedAuthServiceServer
}

func NewGrpcAuthServer(ucase domain.AuthUsecase) *grpcAuthServer {
	return &grpcAuthServer{
		authUsecae: ucase,
	}
}

func (g *grpcAuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user domain.User
	user.Email = req.Email
	user.Password = req.Password

	res := g.authUsecae.Register(ctx, &user)

	return &pb.RegisterResponse{
		Status: res.Status,
		Error:  res.Error,
	}, nil
}

func (g *grpcAuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user domain.User
	user.Email = req.Email
	user.Password = req.Password

	res := g.authUsecae.Login(ctx, &user)
	return &pb.LoginResponse{
		Status: res.Status,
		Error:  res.Error,
		Token:  res.Token,
	}, nil
}

func (g *grpcAuthServer) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	var user domain.User
	user.Token = req.Token

	res := g.authUsecae.Validate(ctx, &user)

	return &pb.ValidateResponse{
		Status: res.Status,
		Error:  res.Error,
		UserId: res.UserId,
	}, nil

}
