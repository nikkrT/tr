package applicaition

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"micr_course/orderService/config"
	usecase "micr_course/orderService/order/business"
	grpc_product "micr_course/orderService/order/infrastructure/grpc"
	postgres "micr_course/orderService/order/infrastructure/postgres"
	"micr_course/orderService/order/infrastructure/postgres/repo"
	myGrpc "micr_course/orderService/order/interfaces"
	pb "micr_course/pkg/proto/orderService"
	pb_product "micr_course/pkg/proto/productService"
	"net"
)

type App struct {
	psql   *pgxpool.Pool
	config *config.Config
}

func NewApp(config *config.Config) *App {
	app := &App{
		config: config,
		psql:   nil,
	}
	return app
}

func (app *App) Start(ctx context.Context) error {

	pool, err := postgres.InitDB(ctx, app.config.Postgres)
	if err != nil {
		return fmt.Errorf("failed to init postgres: %v", err)
	}

	app.psql = pool

	defer app.psql.Close()

	orderRepo := repo.NewProductRepo(pool)
	conn, err := grpc.NewClient(app.config.GRPC.AddressGrpcProduct,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to create grpc connection: %v", err)
	}
	client := pb_product.NewProductServiceClient(conn)
	grpcProduct := grpc_product.NewGRPC(client)

	service := usecase.NewService(orderRepo, grpcProduct)
	grpcInterface := myGrpc.NewGRPCInterface(service)

	ch := make(chan error, 2)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", app.config.GRPC.Port))
		if err != nil {
			ch <- fmt.Errorf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		reflection.Register(s)
		pb.RegisterOrderServiceServer(s, grpcInterface)
		ch <- s.Serve(lis)
		close(ch)
	}()
	select {
	case err := <-ch:
		fmt.Println(err)
		return err
	}
}
