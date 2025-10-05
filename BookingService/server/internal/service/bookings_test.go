package service

import (
	"context"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"

	bookingsV1 "github.com/LashkaPashka/BookingService/bookings_proto/gen/go/booking"
	"github.com/LashkaPashka/BookingService/server/grpcclient"
	"github.com/LashkaPashka/BookingService/server/internal/config"
	convertbody "github.com/LashkaPashka/BookingService/server/internal/lib/convertBody"
	"github.com/LashkaPashka/BookingService/server/internal/storage/postgresql"
)
var in = &bookingsV1.CreateBookingRequest{
	UserId: "507783c0-6a18-40aa-ba61-645500c46e86",
	HotelId: "11111111-1111-1111-1111-111111111111",
	RoomTypeId: "aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
	CheckIn: "2025-09-28",
	CheckOut: "2025-10-05",
	RoomsCount: 4,
	IdempotencyKey: "12acf5be-58d8-4558-a5fd-987f283b573a",
	Payment: &bookingsV1.PaymentInfo{
		Method: "card",
		Token: "tol_test_123",
		TotalAmount: 10000.0,
		Currency: "RUB",
	},
}
var service *BookingService

func TestMain(m *testing.M) {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	clienterGetter := grpcclient.New(logger)

	storage := postgresql.New(*cfg, logger)

	service = New(cfg, logger, storage, storage, clienterGetter)

	m.Run()
}

func TestCreateBooking(t *testing.T) {
	payloadBooking := convertbody.ConvertBody(in)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	booking_id, status := service.CreateBooking(ctx, *payloadBooking)

	if len(booking_id) == 0 {
		t.Fatalf("Length booking_id equal zero")
	}

	want := "PENDING"
	if !reflect.DeepEqual(status, want) {
		t.Fatalf("status not want value")
	}

	t.Log(booking_id, status)
}



func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}