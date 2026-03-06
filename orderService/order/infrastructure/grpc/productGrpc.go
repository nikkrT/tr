package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "micr_course/pkg/proto/productService"
)

type GRPC struct {
	client pb.ProductServiceClient
}

func NewGRPC(client pb.ProductServiceClient) *GRPC {
	return &GRPC{client: client}
}

func (app *GRPC) CheckProductAvailability(ctx context.Context, productId int) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	feature, err := app.client.ReadProduct(ctx, &pb.ReadProductRequest{Id: int64(productId)})
	if err != nil {
		return -1, status.Error(codes.Internal, err.Error())
	}
	return int(feature.Product.Id), nil
}
