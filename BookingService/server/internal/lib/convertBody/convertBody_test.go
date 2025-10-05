package convertbody

import (
	"reflect"
	"testing"

	bookingsV1 "github.com/LashkaPashka/BookingService/bookings_proto/gen/go/booking"
	"github.com/LashkaPashka/BookingService/server/internal/model"
)

const idempotencyKey string = "550e8400-e29b-41d4-a716-446655440000"

var in = &bookingsV1.CreateBookingRequest{
	UserId:         "u123",
	HotelId:        "h456",
	RoomTypeId:     "r1",
	CheckIn:        "2025-10-01",
	CheckOut:       "2025-10-05",
	RoomsCount:     2,
	IdempotencyKey: idempotencyKey,
	Payment: &bookingsV1.PaymentInfo{
		Method:      "card",
		Token:       "tol_test_123",
		TotalAmount: 10000.0,
		Currency:    "RUB",
	},
}

func TestConvertBody(t *testing.T) {
	nightsGot := ParseDate(in.CheckIn, in.CheckOut)

	nightsWant := 4

	if !reflect.DeepEqual(nightsGot, nightsWant) {
		t.Fatalf("not equal got: %d && want:%d", nightsGot, nightsWant)
	}

	gotBookings := ConvertBody(in)

	wantBookings := &model.Booking{
		UserID: in.UserId,
		HotelID: in.HotelId,
		RoomTypeID: in.RoomTypeId,
		CheckIn: in.CheckIn,
		CheckOut: in.CheckOut,
		Nights: nightsGot,
		Status: statusPending,
		TotalAmount: in.Payment.TotalAmount,
		Currency: in.Payment.Currency,
		IdempotencyKey: idempotencyKey,
	}

	if !reflect.DeepEqual(gotBookings, wantBookings) {
		t.Fatalf("not equal got: %v && want:%v", gotBookings, wantBookings)
	}
}
