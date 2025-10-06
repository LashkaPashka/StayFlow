package postgresql

import (
	"context"
	"log"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"

	bookingsV1 "github.com/LashkaPashka/BookingService/bookings_proto/gen/go/booking"
	convertbody "github.com/LashkaPashka/BookingService/server/internal/lib/convertBody"
	"github.com/LashkaPashka/BookingService/server/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	userID         = uuid.New().String()
	hotelID        = uuid.New().String()
	roomTypeId     = uuid.New().String()
	idempotencyKey = uuid.New().String()
)
var st *Storage

var in = &bookingsV1.CreateBookingRequest{
	UserId: userID,
	HotelId: hotelID,
	RoomTypeId: roomTypeId,
	CheckIn: "2025-09-28",
	CheckOut: "2025-10-05",
	RoomsCount: 2,
	IdempotencyKey: idempotencyKey,
	Payment: &bookingsV1.PaymentInfo{
		Method: "card",
		Token: "tol_test_123",
		TotalAmount: 10000.0,
		Currency: "RUB",
	},
}

func TestMain(m *testing.M) {
	logger := setupLogger("local")

	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, "postgres://postgres:root@localhost:5432/BookingsDb")
	if err != nil {
		log.Fatalf("Invalid connection to Db: %s", err.Error())
		os.Exit(1)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Invalid ping to BookingsDb")
		os.Exit(1)
	}

	st = &Storage{
		Pool: pool,
		Logger: logger,
	}

	m.Run()
}


func TestCreateBooking(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bookings := convertbody.ConvertBody(in)

	booking_id, err := st.SaveBooking(ctx, bookings)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	log.Println(booking_id)
}

func TestUpdateBooking(t *testing.T) {
	const status = "CONFIRMED"
	const booking_id = "38f00081-1326-479c-a34d-26726d64c9e8"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gotBooking, err := st.UpdateStatusOfBooking(ctx, booking_id, status)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	wantBookings := model.Booking{
		UserID: userID,
		HotelID: hotelID,
		RoomTypeID: roomTypeId,
		CheckIn: in.CheckIn,
		CheckOut: in.CheckOut,
		Nights: 4,
		RoomsCount: int(in.RoomsCount),
		Status: status,
		TotalAmount: in.Payment.TotalAmount,
		Currency: in.Payment.Currency,
		IdempotencyKey: idempotencyKey,
	}

	if reflect.DeepEqual(gotBooking, wantBookings) {
		t.Fatalf("Not equal - got: %v && want: %v", gotBooking, wantBookings)
	}
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