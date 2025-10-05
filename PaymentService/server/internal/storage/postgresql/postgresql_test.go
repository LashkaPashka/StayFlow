package postgresql

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/model"
)

const (
	userID = "e4a7f8f9-2c53-4d1c-b2d4-4a9395b81c31"
	bookingID = "7b2d9f4a-8c3e-4f6d-9a1b-2e5c7d8f0a12"
)

var payment = &model.Payment{
	UserID: userID,
	BookingID: bookingID,
	Amount: 5000,
	Currency: "RUB",
	Method: "card",
	PaymentID: "12d3d3d3",
	ClientSecret: "d3f23f3",
	Status: "PENDING",
	Response: "WD1D",
} 

var connStr = "postgres://postgres:root@localhost:5432/PaymentDb"
var st *Storage 

func TestMain(m *testing.M) {
	logger := setupLogger("test")

	st = New(connStr, logger)

	m.Run()
}

func TestStoreCreatedPayment(t *testing.T) {
	if err := st.StoreCreatedPayment(context.Background(), payment); err != nil {
		t.Fatal("Test #1: error - StoreCreatedPayment")
	}
	
	t.Log("Ok!")
}

func TestUpdatePayment(t *testing.T) {
	payment.Status = "PENDING"
	
	if err := st.UpdatePayment(context.Background(), payment); err != nil {
		t.Fatal("Test #2: error - UpdatePayment")
	}

	t.Log("Ok!")
}

func TestGetActivePayment(t *testing.T) {
	payment, err := st.GetActivePayment(context.Background(), userID)
	if err != nil {
		t.Fatal("Test #3: error - GetActivePayment")
	}

	if len(payment.UserID) > 0 {
		t.Log("Test #3: payment exist")
	}

	t.Log("Ok!")
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