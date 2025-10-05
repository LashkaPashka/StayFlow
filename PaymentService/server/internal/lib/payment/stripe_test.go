package payment

import (
	"log/slog"
	"os"
	"reflect"

	//"reflect"
	"testing"

	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/config"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/model"
	"github.com/stripe/stripe-go/v78"
	//"github.com/stripe/stripe-go/v78"
)

const (
	userID = "e4a7f8f9-2c53-4d1c-b2d4-4a9395b81c31"
	bookingID = "7b2d9f4a-8c3e-4f6d-9a1b-2e5c7d8f0a12"
	pID = "cs_test_a1xY0dqUVKSlbxvmjnJHRW2ndNc3wG9cftw8moj1ddOw7LASws8K7rLTKr"
)

var payment = &model.Payment{
	UserID: userID,
	BookingID: bookingID,
	Amount: 5000,
	Currency: "RUB",
	Method: "card",
}

var pt *Payment 

func TestMain(m *testing.M) {
	cfg := config.MustLoad()
	
	logger := setupLogger("test")
	pt = New(cfg, logger)

	m.Run()
}

func TestCreatePayment(t *testing.T) {
	sessionResult, err := pt.CreatePayment(payment)
	if err != nil {
		t.Fatal("Test #1: error CreatePayment")
	}

	userID := sessionResult.Metadata["user_id"]
	bookingID := sessionResult.Metadata["booking_id"]
	
	if userID == "" || bookingID == "" {
		t.Fatal("user_id is empty && booking_id is empty")
	}

	t.Log(sessionResult.URL)
}

func TestGetPaymentStatus(t *testing.T) {
	session, err := pt.GetPaymentStatus(pID)
	if err != nil {
		t.Fatal("Test #2: error - GetPaymentStatus")
	}

	gotStatus := session.Status
	wantStatus := stripe.CheckoutSessionStatus("complete")
	if !reflect.DeepEqual(gotStatus, wantStatus) {
		t.Fatalf("Test #2: error - gotStatus: %v != wantStatus: %v", gotStatus, wantStatus)
	}

	t.Log(gotStatus)

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