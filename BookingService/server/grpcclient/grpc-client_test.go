package grpcclient

import (
	"log/slog"
	"os"
	"reflect"
	"testing"

	"github.com/LashkaPashka/BookingService/server/internal/model"
)

const (
	paymentID = "cs_test_a1x2cEZRBMzLhZ5GWogIXoIW7Nk9yyt8ro5xlgQmvJMsxJUD7Q4YtVx6Mj"
	bookingID = "7b2d9f4a-8c3e-4f6d-9a1b-2e5c7d8f0a12"
)
var payloadBookings = model.Booking{
	UserID: "507783c0-6a18-40aa-ba61-645500c46e86",
	HotelID: "11111111-1111-1111-1111-111111111111",
	RoomTypeID: "aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
	CheckIn: "2025-09-28",
	CheckOut: "2025-10-05",
	Currency: "RUB",
	TotalAmount: 5000,
	RoomsCount: 4,
	IdempotencyKey: "e4a7f8f9-2c53-4d1c-b2d4-4a9395b81c31",
}

var grpcClient *Client

func TestMain(m *testing.M) {
	logger := setupLogger("test")

	grpcClient = New(logger)

	m.Run()
}

func TestCheckRoomAvailability(t *testing.T) {
	addr := "localhost:8081"

	available, availableRooms := grpcClient.CheckRoomAvailability(addr, payloadBookings)
	if !available {
		t.Fatalf("Available equal: %v", available)
	}

	wantRooms := 5

	if wantRooms > availableRooms {
		t.Fatalf("Not available rooms")
	}

	t.Log(available, availableRooms)
}

func TestReserveRoom(t *testing.T) {
	addr := "localhost:8081"

	success, err := grpcClient.ReserveRoom(addr, payloadBookings)
	if err != nil {
		t.Fatalf("Error in ReserveRoom: err %v", err.Error())
	}
	
	if reflect.DeepEqual(success, false) {
		t.Fatal("success equal false")
	}
}

func TestReleaseRoom(t *testing.T) {
	addr := "localhost:8081"

	success, err := grpcClient.ReleaseRoom(addr, payloadBookings)
	if err != nil {
		t.Fatalf("Error in ReleaseRoom: err %v", err.Error())
	}
	
	if reflect.DeepEqual(success, false) {
		t.Fatal("success equal false")
	}
	
}

func TestCreatePayment(t *testing.T) {
	addr := "localhost:8082"

	url, paymentID, status := grpcClient.CreatePayment(addr, model.PaymentInfo{
		TotalAmount: payloadBookings.TotalAmount,
		BookingID: bookingID,
		UserID: payloadBookings.UserID,
		Currency: payloadBookings.Currency,
		Method: "card",
		Token: "tok_1QZ3Lh2eZvKYlo2CwZz9V3nb",
		IdempotencyKey: payloadBookings.IdempotencyKey,
	})

	if len(url) == 0 || len(paymentID) == 0 || len(status) == 0 {
		t.Fatal("Length url || paymentID || status equal zero")
	}

	t.Log(url)
	t.Log(status)
}

func TestGetPaymentStatus(t *testing.T) {
	addr := "localhost:8082"

	gotStatus := grpcClient.GetPaymentStatus(addr, paymentID)
	wantStatus := "complete"

	if !reflect.DeepEqual(gotStatus, wantStatus) {
		t.Fatalf("gotStatus: %s != wantStatus: %s", gotStatus, wantStatus)
	}

	t.Log(gotStatus)
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