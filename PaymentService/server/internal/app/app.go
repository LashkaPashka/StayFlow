package app

import (
	"log/slog"

	grpcapp "github.com/LashkaPashka/StayFlow/PaymentService/server/internal/app/grpc"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/config"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	grpcPort string,
	storagePath string,
) *App {
	grpcApp := grpcapp.New(cfg, logger, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
