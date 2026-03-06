package business

import (
	"context"
	"fmt"
	"micr_course/orderService/paymentService"
	"micr_course/pkg/models"
)

type OrderRepository interface {
	Create(ctx context.Context, order models.OrderModel) (int, error)
	ReadById(ctx context.Context, id int) (models.OrderModel, error)
	Update(ctx context.Context, order models.OrderModel) error
	DeleteById(ctx context.Context, id int) error
}

type GRPC interface {
	CheckProductAvailability(ctx context.Context, productId int) (int, error)
}

type UseCase struct {
	orderRepository OrderRepository
	grpcProduct     GRPC
}

func NewService(repo OrderRepository, grpc GRPC) *UseCase {
	return &UseCase{
		orderRepository: repo,
		grpcProduct:     grpc,
	}
}

func (s *UseCase) UseCase(ctx context.Context, productId int) (int, error) {

	res, err := s.grpcProduct.CheckProductAvailability(ctx, productId)
	if err != nil {
		return -1, fmt.Errorf("grpc product availability error: %w", err)
	}
	if res <= 0 {
		fmt.Println("продукт не существует в бд")
		return -1, nil
	}
	fmt.Printf("продукт существует в бд с номером %d", res)
	if paymentService.CheckPayment(productId) == false {
		fmt.Println("оплата неудалась к сожалению")
		return -1, nil
	}
	order := models.OrderModel{
		ProductId: productId,
		Status:    "created",
	}
	orderId, err := s.orderRepository.Create(ctx, order)
	fmt.Println("заказ успешно создан с Id=%d", orderId)
	return orderId, nil
}
