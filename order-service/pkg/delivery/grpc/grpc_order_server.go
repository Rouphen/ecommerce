package orderserver

import (
	"context"
	"net/http"

	"ecommerce/order-service/pkg/delivery/grpc/pb"
	domain "ecommerce/order-service/pkg/domain"
)

type grpcOrderServer struct {
	ucase        domain.OrderUsecase
	productClent *grpcProductClient

	pb.UnimplementedOrderServiceServer
}

func NewGrpcOrderServer(ucase domain.OrderUsecase, productClent *grpcProductClient) *grpcOrderServer {
	return &grpcOrderServer{
		ucase:        ucase,
		productClent: productClent,
	}
}

func (s *grpcOrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	product, err := s.productClent.FindOne(req.ProductId)

	if err != nil {
		return &pb.CreateOrderResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	}

	if product.Status >= http.StatusNotFound {
		return &pb.CreateOrderResponse{Status: product.Status, Error: product.Error}, nil
	}

	if product.Data.Stock < req.Quantity {
		return &pb.CreateOrderResponse{Status: http.StatusConflict, Error: "Stock too less"}, nil
	}

	order := domain.Order{
		Price:     product.Data.Price,
		ProductId: product.Data.Id,
		UserId:    req.UserId,
	}

	s.ucase.Create(ctx, &order)

	res, err := s.productClent.DecreaseStock(req.ProductId, order.Id)

	if err != nil {
		return &pb.CreateOrderResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	}

	if res.Status == http.StatusConflict {
		s.ucase.Delete(ctx, order.Id)
		return &pb.CreateOrderResponse{Status: http.StatusConflict, Error: res.Error}, nil
	}

	return &pb.CreateOrderResponse{
		Status: http.StatusCreated,
		Id:     order.Id,
	}, nil
}
