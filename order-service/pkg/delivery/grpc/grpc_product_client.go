package orderserver

import (
	"context"

	"ecommerce/order-service/pkg/delivery/grpc/pb"
	"google.golang.org/grpc"
)

type grpcProductClient struct {
	client pb.ProductServiceClient
}

func NewGrpcProductClient(cc grpc.ClientConnInterface) *grpcProductClient {
	return &grpcProductClient{
		client: pb.NewProductServiceClient(cc),
	}
}

func (g *grpcProductClient) FindOne(productId int64) (*pb.FindOneResponse, error) {
	req := &pb.FindOneRequest{
		Id: productId,
	}

	return g.client.FindOne(context.Background(), req)
}

func (g *grpcProductClient) DecreaseStock(productId int64, orderId int64) (*pb.DecreaseStockResponse, error) {
	req := &pb.DecreaseStockRequest{
		Id:      productId,
		OrderId: orderId,
	}

	return g.client.DecreaseStock(context.Background(), req)
}
