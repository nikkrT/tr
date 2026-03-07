package interfaces

import (
	"context"
	model "micr_course/pkg/models"
	"micr_course/pkg/proto/orderService"
	pb "micr_course/pkg/proto/productService"
)

type Service interface {
	UseCase(ctx context.Context, productId int) (int, error)
}

type GRPCServer struct {
	orderService.UnimplementedOrderServiceServer
	service Service
}

func NewGRPCInterface(service Service) *GRPCServer {
	return &GRPCServer{
		service: service,
	}
}

func (g *GRPCServer) CreateOrder(ctx context.Context, req *orderService.OrderRequest) (*orderService.OrderAnswer, error) {
	input := int(req.GetId())
	res, err := g.service.UseCase(ctx, input)
	if err != nil {
		return &orderService.OrderAnswer{}, err
	}
	if res > 0 {
		return &orderService.OrderAnswer{OrderId: uint32(res), Success: true}, nil
	}
	return &orderService.OrderAnswer{}, nil
}

func convertGRPCRequestToProduct(req *pb.ReadProductResponse) model.Product {
	return model.Product{
		Id:          int(req.Product.Id),
		Name:        req.Product.Name,
		Price:       int(req.Product.Price),
		Description: req.Product.Description,
	}
}
