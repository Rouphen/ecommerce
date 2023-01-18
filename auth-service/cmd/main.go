package main

import (
	"fmt"
	"log"
	"net"

	"ecommerce/auth-service/pkg/config"
	grpchandler "ecommerce/auth-service/pkg/delivery/grpc"
	"ecommerce/auth-service/pkg/delivery/grpc/pb"
	repo "ecommerce/auth-service/pkg/repository"
	ucase "ecommerce/auth-service/pkg/usecase"
	"ecommerce/auth-service/pkg/utils"

	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	db, _ := utils.ConnectMySQLByGorm(c.DBUrl)
	authRepo := repo.NewUserRepository(db)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "auth-service",
		ExpirationHours: 24 * 365,
	}
	authUsecase := ucase.NewAuthUsecase(authRepo, jwt)
	handler := grpchandler.NewGrpcAuthServer(authUsecase)

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Auth Svc on", c.Port)

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, handler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
