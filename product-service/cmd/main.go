package main

import (
	"fmt"
	"log"
	"net"

	"ecommerce/product-service/pkg/config"
	grpchandler "ecommerce/product-service/pkg/delivery/grpc"
	pb "ecommerce/product-service/pkg/delivery/grpc/pb"
	domain "ecommerce/product-service/pkg/domain"
	repo "ecommerce/product-service/pkg/repository"
	ucase "ecommerce/product-service/pkg/usecase"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	c := config.NewConfig().InLocalConfig().InDockerComposeEnv()
	db, _ := ConnectMySQLByGorm(c.DBUrl)

	productRepo := repo.NewProdcutRepository(db)
	stocklogRepo := repo.NewStockDecreaseLogsRepository(db)

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}
	fmt.Println("Product Svc on", c.Port)

	productUsecase := ucase.NewProductUsecase(productRepo, stocklogRepo)
	service := grpchandler.NewGrpcProductService(productUsecase)
	grpcServer := grpc.NewServer()

	pb.RegisterProductServiceServer(grpcServer, service)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}

func ConnectMySQLByGorm(url string) (*gorm.DB, error) {
	log.Println(url)
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&domain.Product{})
	db.AutoMigrate(&domain.StockDecreaseLog{})

	return db, nil
}
