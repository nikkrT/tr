package grpc

import (
	"context"
	"micr_course/pkg/models"
	"micr_course/pkg/proto/productService"
	"micr_course/productService/product/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler interface {
	CreateProduct(ctx context.Context, product models.Product) error
	ReadProduct(ctx context.Context, id int) (models.Product, error)
	UpdateProduct(ctx context.Context, product models.Product) error
	DeleteProduct(ctx context.Context, id int) error
}

type GRPCServer struct {
	productService.UnimplementedProductServiceServer
	service GrpcHandler
}

func NewGRPCServer(service *service.ProductService) *GRPCServer {
	return &GRPCServer{service: service}
}

func (g *GRPCServer) CreateProduct(ctx context.Context, req *productService.CreateProductRequest) (*productService.CreateProductResponse, error) {
	input := models.Product{
		Name:        req.GetName(),
		Price:       int(req.GetPrice()),
		Description: req.GetDescription(),
	}
	err := g.service.CreateProduct(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &productService.CreateProductResponse{Id: 1}, nil
}

func (g *GRPCServer) ReadProduct(ctx context.Context, req *productService.ReadProductRequest) (*productService.ReadProductResponse, error) {
	input := int(req.GetId())
	output, err := g.service.ReadProduct(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &productService.ReadProductResponse{Product: convertProductToPb(output)}, nil
}

func (g *GRPCServer) UpdateProduct(ctx context.Context, req *productService.UpdateProductRequest) (*productService.UpdateProductResponse, error) {
	input := models.Product{
		Name:        req.GetName(),
		Price:       int(req.GetPrice()),
		Description: req.GetDescription(),
	}
	err := g.service.UpdateProduct(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &productService.UpdateProductResponse{Result: true}, nil
}

func (g *GRPCServer) DeleteProduct(ctx context.Context, req *productService.DeleteProductRequest) (*productService.DeleteProduceResponse, error) {
	input := int(req.GetId())
	err := g.service.DeleteProduct(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &productService.DeleteProduceResponse{Result: true}, nil
}

func convertProductToPb(product models.Product) *productService.Product {
	return &productService.Product{
		Id:          int64(product.Id), // приводим int к int64
		Name:        product.Name,
		Description: product.Description,
		Price:       uint32(product.Price),
	}
}
