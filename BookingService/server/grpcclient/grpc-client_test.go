package grpcclient

import (
	"log/slog"
	"os"
	"reflect"
	"testing"

	"github.com/LashkaPashka/BookingService/server/internal/model"
)

var payloadBookings = model.Booking{
	UserID: "us_01f3gf2",
	HotelID: "11111111-1111-1111-1111-111111111111",
	RoomTypeID: "aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
	CheckIn: "2025-09-28",
	CheckOut: "2025-10-05",
	RoomsCount: 4,
}

var grpcClient *Client

func TestMain(m *testing.M) {
	logger := setupLogger("local")

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