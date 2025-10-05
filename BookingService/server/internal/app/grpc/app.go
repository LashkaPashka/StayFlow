package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/LashkaPashka/BookingService/server/grpcclient"
	"github.com/LashkaPashka/BookingService/server/internal/config"
	bookingGRPC "github.com/LashkaPashka/BookingService/server/internal/grpc/booking"
	"github.com/LashkaPashka/BookingService/server/internal/service"
	"github.com/LashkaPashka/BookingService/server/internal/storage/postgresql"
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

	gRPCServer := grpc.NewServer()

	storage := postgresql.New(*cfg, logger)

	clientGRPC := grpcclient.New(logger)

	bookingService := service.New(cfg, logger, storage, storage, clientGRPC)

	bookingGRPC.Register(gRPCServer, bookingService)

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

func (a *App) Run() error {
	const op = "gprcapp.Run"

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

	return  nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.logger.With(slog.String("op", op)).
			Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}