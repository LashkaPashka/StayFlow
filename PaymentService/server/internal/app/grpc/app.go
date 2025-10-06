package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/config"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/grpc/payment"
	paymentLib "github.com/LashkaPashka/StayFlow/PaymentService/server/internal/lib/payment"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/service"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/storage/postgresql"
	"google.golang.org/grpc"
)

const host = "0.0.0.0"

type App struct {
	logger     *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	port string,
) *App {

	gRPCServer := grpc.NewServer()

	storage := postgresql.New(cfg.StoragePath, logger)

	paymentClient := paymentLib.New(cfg, logger)

	service := service.New(storage, paymentClient, cfg, logger)

	payment.Register(service, gRPCServer)

	return &App{
		logger:     logger,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.logger.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)

	l, err := net.Listen("tcp", net.JoinHostPort(host, a.port))
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.logger.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.String("port", a.port))

	a.gRPCServer.GracefulStop()
}
