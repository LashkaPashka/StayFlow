package service

import (
	"context"
	"log/slog"
	"os"
	"testing"

	paymentV1 "github.com/LashkaPashka/StayFlow/PaymentService/payment_proto/gen/go/payment"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/config"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/lib/payment"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/storage/postgresql"
)

const (
	userID    = "e4a7f8f9-2c53-4d1c-b2d4-4a9395b81c31"
	bookingID = "7b2d9f4a-8c3e-4f6d-9a1b-2e5c7d8f0a12"
)

var payload = &paymentV1.CreatePaymentRequest{
	BookingId:      bookingID,
	UserId:         userID,
	Amount:         5000,
	Currency:       "RUB",
	Method:         "card",
	Token:          "tok_1QZ3Lh2eZvKYlo2CwZz9V3nb",
	IdempotencyKey: "",
}

var sv *Service

func TestMain(m *testing.M) {
	logger := setupLogger("test")
	cfg := config.MustLoad()
	storage := postgresql.New(cfg.StoragePath, logger)
	paymentClient := payment.New(cfg, logger)

	sv = New(storage, paymentClient, cfg, logger)

	m.Run()
}

func TestCreatePayment(t *testing.T) {
	_ = payload

	sv.CreatePayment(context.Background(), nil)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "test":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
