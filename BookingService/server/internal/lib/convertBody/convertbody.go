package convertbody

import (
	"time"

	bookingsV1 "github.com/LashkaPashka/BookingService/bookings_proto/gen/go/booking"
	//"github.com/LashkaPashka/BookingService/server/internal/lib/random"
	"github.com/LashkaPashka/BookingService/server/internal/model"
)

// var testCreateBooking = &bookingsV1.CreateBookingRequest{
// 	UserId: "u123",
// 	HotelId: "h456",
// 	RoomTypeId: "r1",
// 	CheckIn: "2025-10-01",
// 	CheckOut: "2025-10-05",
// 	RoomsCount: 2,
// 	IdempotencyKey: "550e8400-e29b-41d4-a716-446655440000",
// 	Payment: &bookingsV1.PaymentInfo{
// 		Method: "card",
// 		Token: "tol_test_123",
// 		TotalAmount: 10000.0,
// 		Currency: "RUB",
// 	},
// }



const (
	statusPending = "PENDING"
)

func ConvertBody(in *bookingsV1.CreateBookingRequest) model.Booking {
	nights := ParseDate(in.CheckIn, in.CheckOut)

	return model.Booking{
		UserID: in.UserId,
		HotelID: in.HotelId,
		RoomTypeID: in.RoomTypeId,
		CheckIn: in.CheckIn,
		CheckOut: in.CheckOut,
		Nights: nights,
		RoomsCount: int(in.RoomsCount),
		Status: statusPending,
		TotalAmount: in.Payment.TotalAmount,
		Currency: in.Payment.Currency,
		IdempotencyKey: in.IdempotencyKey,
	}
}

func ParseDate(checkIn, checkOut string) int {
	layout := "2006-01-02"
	
	parseCheckIn, _ := time.Parse(layout, checkIn)
	parseCheckOut, _ := time.Parse(layout, checkOut)

	dateCheckIn := time.Date(parseCheckIn.Year(), parseCheckIn.Month(), parseCheckIn.Day(), 0, 0, 0, 0, time.UTC)
	dateCheckOut := time.Date(parseCheckOut.Year(), parseCheckOut.Month(), parseCheckOut.Day(), 0, 0, 0, 0, time.UTC)

	diff := dateCheckOut.Sub(dateCheckIn)

	return int(diff.Hours()/24)
}