package postgresql

import (
	"context"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"
)

type TestHotel struct {
	HotelID    string `json:"hotel_id"`
	RoomTypeID string `json:"room_type_id"`
	CheckIn    string `json:"check_in"`
	CheckOut   string `json:"check_out"`
	RoomsCount int    `json:"rooms_count"`
}

var payload = &TestHotel{
	HotelID:    "11111111-1111-1111-1111-111111111111",
	RoomTypeID: "aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
	CheckIn:    "2025-09-28",
	CheckOut:   "2025-10-05",
	RoomsCount: 5,
}

const connStr = "postgres://postgres:root@localhost:5432/postgres"

var st *Storage

func TestMain(m *testing.M) {
	logger := setupLogger("local")

	st = New(connStr, logger)

	m.Run()
}

func TestGetAvailableRooms(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	availableRooms, err := st.GetAvailableRooms(
		ctx,
		payload.HotelID,
		payload.RoomTypeID,
		payload.CheckIn,
		payload.CheckOut,
	)
	if err != nil {
		t.Fatalf("Test №1 - error: %v", err.Error())
	}

	if payload.RoomsCount > availableRooms {
		t.Fatalf("Not available rooms")
	}

	t.Logf("Quantity available rooms equal %d", availableRooms)
}

func TestBookedRooms(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	success, err := st.BookedRoom(
		ctx,
		payload.HotelID,
		payload.RoomTypeID,
		payload.CheckIn,
		payload.CheckOut,
		payload.RoomsCount,
	)
	if err != nil {
		t.Fatalf("Test №2 - error: %v", err)
	}

	if !reflect.DeepEqual(success, true) {
		t.Fatalf("Not equal")
	}
}

func TestFreeRoom(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	success, err := st.FreeRoom(
		ctx,
		payload.HotelID,
		payload.RoomTypeID,
		payload.CheckIn,
		payload.CheckOut,
		payload.RoomsCount,
	)

	if err != nil {
		t.Fatalf("Test №2 - error: %v", err)
	}

	if !reflect.DeepEqual(success, true) {
		t.Fatalf("Not equal")
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
