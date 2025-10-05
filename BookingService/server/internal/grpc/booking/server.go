package booking

import (
	"context"
	"errors"
	"strings"

	bookingsV1 "github.com/LashkaPashka/BookingService/bookings_proto/gen/go/booking"
	convertbody "github.com/LashkaPashka/BookingService/server/internal/lib/convertBody"
	"github.com/LashkaPashka/BookingService/server/internal/model"
	"google.golang.org/grpc"
)

type BookingService interface {
	CreateBooking(ctx context.Context, booking model.Booking) (booking_id string, status string)
	ConfirmBooking(ctx context.Context, booking_id, payment_id string) model.Booking
	CancelBooking(ctx context.Context, booking_id, paymendID, user_id, reason string) string
	GetBooking(ctx context.Context, booking_id, user_id string) model.Booking
	ListBookings(ctx context.Context, user_id string) []model.Booking
}

type serverAPI struct {
	bookingsV1.UnimplementedBookingsServiceServer
	bookingService BookingService
}

func Register(gRPCServer *grpc.Server, bookingService BookingService) {
	bookingsV1.RegisterBookingsServiceServer(gRPCServer, &serverAPI{bookingService: bookingService})
}

func (s *serverAPI) CreateBooking(
	ctx context.Context, 
	in *bookingsV1.CreateBookingRequest,
) (*bookingsV1.CreateBookingsResponse, error) {
	// TODO: convert Body in model for postgreSQL
	booking := convertbody.ConvertBody(in)

	// TODO: call CreateBooking from bookingService
	booking_id, status := s.bookingService.CreateBooking(ctx, *booking)
	if strings.Compare(booking_id, status) == 0 {
		return nil, errors.New("error inside bookingService")
	}

	return &bookingsV1.CreateBookingsResponse{
		BookingId: booking_id,
		Status: status,
	}, nil
}

func (s *serverAPI) ConfirmBooking(
	ctx context.Context, 
	in *bookingsV1.ConfirmBookingRequest,
) (*bookingsV1.Booking, error) {

	updatedBooking := s.bookingService.ConfirmBooking(ctx, in.BookingId, in.PaymentId)

	return &bookingsV1.Booking{
		BookingId: "",
		UserId: updatedBooking.UserID,
		HotelId: updatedBooking.HotelID,
		RoomTypeId: updatedBooking.RoomTypeID,
		Status: updatedBooking.Status,
		CheckIn: updatedBooking.CheckIn,
		CheckOut: updatedBooking.CheckOut,
		CreatedAt: updatedBooking.CreatedAt.Format("2006-01-02"),
		UpdatedAt: updatedBooking.UpdatedAt.Format("2006-01-02"),
	}, nil
}

func (s *serverAPI) CancelBooking(
	ctx context.Context, 
	in *bookingsV1.CancelBookingRequest,
) (*bookingsV1.Booking, error) {
	

	return &bookingsV1.Booking{
		BookingId: "1",
		Status: "CONFIRMED",
	}, nil
}

func (s *serverAPI) GetBooking(
	ctx context.Context, 
	in *bookingsV1.GetBookingRequest,
) (*bookingsV1.Booking, error) {
	booking := s.bookingService.GetBooking(ctx, in.BookingId, in.UserId)

	return &bookingsV1.Booking{
		BookingId: "",
		UserId: booking.UserID,
		HotelId: booking.HotelID,
		RoomTypeId: booking.RoomTypeID,
		Status: booking.Status,
		CheckIn: booking.CheckIn,
		CheckOut: booking.CheckOut,
		CreatedAt: booking.CreatedAt.Format("2006-01-02"),
		UpdatedAt: booking.UpdatedAt.Format("2006-01-02"),
	}, nil
}

func (s *serverAPI) ListUserBookings(
	ctx context.Context, 
	in *bookingsV1.ListUserBookingsRequest,
) (*bookingsV1.ListUserBookingsReponse, error) {

	return &bookingsV1.ListUserBookingsReponse{
		Bookings: []*bookingsV1.Booking{},
	}, nil
}