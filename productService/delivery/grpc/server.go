package grpc

import (
	"context"
	"productService/internal/model"
	"productService/internal/service"
	pb "productService/pkg/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler interface {
	CreateProduct(ctx context.Context, product model.Product) error
	ReadProduct(ctx context.Context, id int) (model.Product, error)
	UpdateProduct(ctx context.Context, product model.Product) error
	DeleteProduct(ctx context.Context, id int) error
}

type GRPCServer struct {
	pb.UnimplementedProductServiceServer
	service GrpcHandler
}

func NewGRPCServer(service *service.ProductService) *GRPCServer {
	return &GRPCServer{service: service}
}

func (g *GRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	input := model.Product{
		Name:        req.GetName(),
		Price:       int(req.GetPrice()),
		Description: req.GetDescription(),
	}
	err := g.service.CreateProduct(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateProductResponse{Id: 1}, nil
}

func (g *GRPCServer) ReadProduct(ctx context.Context, req *pb.ReadProductRequest) (*pb.ReadProductResponse, error) {
	input := int(req.GetId())
	output, err := g.service.ReadProduct(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ReadProductResponse{Product: convertProductToPb(output)}, nil
}

func (g *GRPCServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	input := model.Product{
		Name:        req.GetName(),
		Price:       int(req.GetPrice()),
		Description: req.GetDescription(),
	}
	err := g.service.UpdateProduct(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UpdateProductResponse{Result: true}, nil
}

func (g *GRPCServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProduceResponse, error) {
	input := int(req.GetId())
	err := g.service.DeleteProduct(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteProduceResponse{Result: true}, nil
}

func convertProductToPb(product model.Product) *pb.Product {
	return &pb.Product{
		Id:          int64(product.Id), // приводим int к int64
		Name:        product.Name,
		Description: product.Description,
		Price:       uint32(product.Price), // приводим int к uint32
		// Если в proto нет полей CreatedAt, UpdatedAt, DeletedAt, их можно пропустить,
		// либо, если они есть и описаны как тип Timestamp, нужно дополнительно конвертировать:
		// CreatedAt: ptypes.TimestampNow(),
	}
}
