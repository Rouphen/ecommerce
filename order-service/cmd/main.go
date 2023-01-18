package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"ecommerce/order-service/pkg/config"
	ogrpc "ecommerce/order-service/pkg/delivery/grpc"
	pb "ecommerce/order-service/pkg/delivery/grpc/pb"
	"ecommerce/order-service/pkg/domain"
	repo "ecommerce/order-service/pkg/repository"
	usecase "ecommerce/order-service/pkg/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	db, _ := ConnectMySQLByGorm(c.DBUrl)
	orderRepo := repo.NewOrderRepository(db)

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	cc, err := grpc.Dial(c.ProductSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(errors.New("Could not connect to Product Server:" + err.Error()))
	}

	productClient := ogrpc.NewGrpcProductClient(cc)

	ucase := usecase.NewOrderUsecae(orderRepo)

	fmt.Println("Order Svc on", c.Port)
	orderService := ogrpc.NewGrpcOrderServer(ucase, productClient)
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, orderService)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}

func ConnectMySQLByGorm(url string) (*gorm.DB, error) {
	log.Println(url)
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
		//return nil, err
	}

	db.AutoMigrate(&domain.Order{})

	return db, nil
}
