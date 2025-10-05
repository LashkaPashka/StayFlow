package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/LashkaPashka/HotelService/server/internal/config"
	hotelgrpc "github.com/LashkaPashka/HotelService/server/internal/grpc/hotel"
	"github.com/LashkaPashka/HotelService/server/internal/services/hotel"
	"github.com/LashkaPashka/HotelService/server/internal/storage/postgresql"
	"google.golang.org/grpc"
)

type App struct {
	logger *slog.Logger
	gRPCServer *grpc.Server
	port int
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	port int,
) *App {

	// TODO: init storage
	storage := postgresql.New(cfg.StoragePath, logger)

	// TODO: init service
	hotelService := hotel.New(storage, cfg, logger)

	gRPCServer := grpc.NewServer()
	hotelgrpc.Register(gRPCServer, hotelService)

	return &App{
		logger: logger,
		gRPCServer: gRPCServer,
		port: port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func(a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.logger.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
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
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}