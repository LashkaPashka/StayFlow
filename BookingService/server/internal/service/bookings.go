package service

import (
	"context"
	"log/slog"

	"github.com/LashkaPashka/BookingService/server/internal/config"
	"github.com/LashkaPashka/BookingService/server/internal/model"
)

const (
	statusPending   = "PENDING"
	statusConfirmed = "CONFIRMED"
)

type Clienter interface {
	CheckRoomAvailability(addr string, booking model.Booking) (available bool, availableRooms int)
	ReserveRoom(addr string, booking model.Booking) (success bool, err error)
	ReleaseRoom(addr string, booking model.Booking) (success bool, err error)
	CreatePayment(addr string, payment model.PaymentInfo) (url, paymentID, status string)
	GetPaymentStatus(addr string, paymendID string) (status string)
	RefundPayment(addr string, paymendID string) (success bool, err error)
}

type BookingSaver interface {
	SaveBooking(
		ctx context.Context,
		booking model.Booking,
	) (booking_id string, err error)
	UpdateStatusOfBooking(
		ctx context.Context,
		booking_id string,
		status string,
	) (booking model.Booking, err error)
}

type BookingGetter interface {
	GetBooking(
		ctx context.Context,
		booking_id string,
		user_id string,
	) (model.Booking, error)
	GetBookings(
		ctx context.Context,
		user_id string,
	) ([]model.Booking, error)
}

type BookingService struct {
	cfg           *config.Config
	logger        *slog.Logger
	bookingSaver  BookingSaver
	bookingGetter BookingGetter
	grpcclienter  Clienter
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	bookingSaver BookingSaver,
	bookingGetter BookingGetter,
	clienter Clienter,
) *BookingService {
	return &BookingService{
		cfg:           cfg,
		logger:        logger,
		bookingSaver:  bookingSaver,
		bookingGetter: bookingGetter,
		grpcclienter:  clienter,
	}
}

func (b *BookingService) CreateBooking(ctx context.Context, booking model.Booking) (url, booking_id, status string) {
	const op = "BookingService.service.bookings.CreateBooking"

	// 1. query on HotelService
	available, avaiableRooms := b.grpcclienter.CheckRoomAvailability(b.cfg.AddrHotelService, booking)

	if !available && avaiableRooms < booking.RoomsCount {
		return "", "", "FAILED"
	}

	// 2. if hotelservice turned true then create booking (STATUS PENDING)
	booking_id, err := b.bookingSaver.SaveBooking(ctx, booking)
	if err != nil {
		b.logger.Error("Invalid save booking in PostgreSQL", slog.String("op", op))
		return "", "", "FAILED"
	}

	// 3. create reserveRoom
	success, err := b.grpcclienter.ReserveRoom(b.cfg.AddrHotelService, booking)
	if err != nil {
		b.logger.Error("Invalid reservedRoom",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)
		return "aaa", "", "FAILED"
	}

	if !success {
		b.logger.Error("Error reserved room", slog.String("op", op))
		return "", "", "FAILED"
	}

	// 4. create payment
	url, _, status = b.grpcclienter.CreatePayment(b.cfg.AddrPaymentService, model.PaymentInfo{
		UserID:         booking.UserID,
		BookingID:      booking_id,
		TotalAmount:    booking.TotalAmount,
		Currency:       booking.Currency,
		Method:         "card",
		Token:          "tok_1QZ3Lh2eZvKYlo2CwZz9V3nb",
		IdempotencyKey: booking.IdempotencyKey,
	})

	return url, booking_id, status
}

func (b *BookingService) ConfirmBooking(ctx context.Context, booking_id, payment_id string) model.Booking {
	const op = "BookingService.service.bookings.ConfirmBooking"

	paymentStatus := b.grpcclienter.GetPaymentStatus(b.cfg.StoragePath, payment_id)

	booking, err := b.bookingSaver.UpdateStatusOfBooking(ctx, booking_id, paymentStatus)
	if err != nil {
		b.logger.Error("Invalid updatedofbookings", slog.String("op", op))
		return model.Booking{}
	}

	return booking
}

func (b *BookingService) CancelBooking(ctx context.Context, booking_id, paymendID, user_id, reason string) string {
	const op = "BookingService.service.bookings.CancelBooking"

	success, err := b.grpcclienter.RefundPayment(b.cfg.AddrPaymentService, paymendID)
	if err != nil {
		b.logger.Error("Invalid reserveRoom",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return ""
	}

	if !success {
		return "failed"
	}

	_, err = b.grpcclienter.ReleaseRoom(b.cfg.AddrHotelService, model.Booking{})
	if err != nil {
		b.logger.Error("Invalid reserverRoom",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)
		return ""
	}

	return ""
}

func (b *BookingService) GetBooking(ctx context.Context, booking_id, user_id string) model.Booking {
	const op = "BookingService.service.bookings.GetBooking"

	getBooking, err := b.bookingGetter.GetBooking(context.Background(), booking_id, user_id)
	if err != nil {
		b.logger.Error("Invalid get booking", slog.String("op", op))
		return model.Booking{}
	}

	return getBooking
}

func (b *BookingService) ListBookings(ctx context.Context, user_id string) []model.Booking {
	const op = "BookingService.service.bookings.ListBooking"

	listBookings, err := b.bookingGetter.GetBookings(context.Background(), user_id)
	if err != nil {
		b.logger.Error("Invalid error getBookings", slog.String("op", op))
		return nil
	}

	return listBookings
}
