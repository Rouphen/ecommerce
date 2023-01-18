package productserver

import (
	"context"
	pb "ecommerce/product-service/pkg/delivery/grpc/pb"
	"ecommerce/product-service/pkg/domain"
)

type grpcProdcutService struct {
	ucase domain.ProductUsecase
	pb.UnimplementedProductServiceServer
}

func NewGrpcProductService(ucase domain.ProductUsecase) *grpcProdcutService {
	return &grpcProdcutService{
		ucase: ucase,
	}
}

func (s *grpcProdcutService) CreateProduct(ctx context.Context,
	req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {

	var product domain.Product

	product.Name = req.Name
	product.Stock = req.Stock
	product.Price = req.Price

	response := s.ucase.Create(ctx, &product)

	return &pb.CreateProductResponse{
		Status: response.Status,
		Error:  response.Error,
		Id:     product.Id,
	}, nil
}

func (s *grpcProdcutService) FindOne(ctx context.Context,
	req *pb.FindOneRequest) (*pb.FindOneResponse, error) {

	product, response := s.ucase.GetByID(ctx, req.Id)
	data := &pb.FindOneData{
		Id:    product.Id,
		Name:  product.Name,
		Stock: product.Stock,
		Price: product.Price,
	}

	return &pb.FindOneResponse{
		Status: response.Status,
		Error:  response.Error,
		Data:   data,
	}, nil
}

func (s *grpcProdcutService) DecreaseStock(ctx context.Context,
	req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {

	response := s.ucase.DecreaseStock(ctx, req.Id, req.OrderId)

	return &pb.DecreaseStockResponse{
		Status: response.Status,
		Error:  response.Error,
	}, nil
}
