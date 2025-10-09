package grpcclient

import (
	"context"
	"log/slog"
	"time"

	hotelsV1 "github.com/LashkaPashka/BookingService/bookings_proto/gen/go/hotel"
	paymentV1 "github.com/LashkaPashka/BookingService/bookings_proto/gen/go/payment"
	"github.com/LashkaPashka/BookingService/server/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	logger *slog.Logger
}

func New(
	logger *slog.Logger,
) *Client {
	return &Client{
		logger: logger,
	}
}

func (c *Client) CheckRoomAvailability(addr string, booking model.Booking) (available bool, availableRooms int) {
	const op = "BookingService.server.grpcclient.CheckRoomAvailability"

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		c.logger.Error("Error create connection")
		return false, -1
	}

	defer conn.Close()

	client := hotelsV1.NewHotelServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.CheckAvailability(ctx, &hotelsV1.CheckAvailabilityRequest{
		HotelId:    booking.HotelID,
		RoomTypeId: booking.RoomTypeID,
		CheckIn:    booking.CheckIn,
		CheckOut:   booking.CheckOut,
	})

	if err != nil {
		c.logger.Error("Invalid grpc query",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)
		return false, -1
	}

	return resp.Available, int(resp.AvailableRooms)
}

func (c *Client) ReserveRoom(addr string, booking model.Booking) (success bool, err error) {
	const op = "BookingService.server.grpcclient.ReserveRoom"

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		c.logger.Error("Error create connection",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)
		return false, err
	}

	defer conn.Close()

	client := hotelsV1.NewHotelServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ReserveRoom(ctx, &hotelsV1.ReserveRoomRequest{
		HotelId:    booking.HotelID,
		RoomTypeId: booking.RoomTypeID,
		CheckIn:    booking.CheckIn,
		CheckOut:   booking.CheckOut,
		RoomsCount: int32(booking.RoomsCount),
	})
	if err != nil {
		c.logger.Error("Invalid grpc query")
		return false, err
	}

	return resp.Success, nil
}

func (c *Client) ReleaseRoom(addr string, booking model.Booking) (success bool, err error) {
	const op = "BookingService.server.grpcclient.ReleaseRoom"

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		c.logger.Error("Error create connection")
		return false, err
	}

	defer conn.Close()

	client := hotelsV1.NewHotelServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ReleaseRoom(ctx, &hotelsV1.ReleaseRoomRequest{
		HotelId:    booking.HotelID,
		RoomTypeId: booking.RoomTypeID,
		CheckIn:    booking.CheckIn,
		CheckOut:   booking.CheckOut,
		RoomsCount: int32(booking.RoomsCount),
	})
	if err != nil {
		c.logger.Error("Invalid grpc query",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)
		return false, err
	}

	return resp.Success, nil
}

func (c *Client) CreatePayment(addr string, payment model.PaymentInfo) (url, paymentID, status string) {
	const op = "BookingService.server.grpcclient.CreatePayment"

	
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		c.logger.Error("Invalid create connection")
		return "", "", ""
	}
	defer conn.Close()

	client := paymentV1.NewPaymentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.CreatePayment(ctx, &paymentV1.CreatePaymentRequest{
		BookingId:      payment.BookingID,
		UserId:         payment.UserID,
		Amount:         payment.TotalAmount,
		Currency:       payment.Currency,
		Method:         payment.Method,
		Token:          payment.Token,
		IdempotencyKey: payment.IdempotencyKey,
	})
	if err != nil {
		c.logger.Error("Invalid grpc query",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return "", "", ""
	}

	return resp.Url, resp.PaymentId, resp.Status
}

func (c *Client) GetPaymentStatus(addr string, paymendID string) (status string) {
	const op = "BookingService.server.grpcclient.GetPaymentStatus"

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		c.logger.Error("Invalid create connection")
		return "Failed"
	}
	defer conn.Close()

	client := paymentV1.NewPaymentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetPaymentStatus(ctx, &paymentV1.GetPaymentStatusRequest{
		PaymentId: paymendID,
	})
	if err != nil {
		c.logger.Error("Invalid grpc query",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return "Failed"
	}

	return resp.Status
}

func (c *Client) RefundPayment(addr string, paymendID string) (success bool, err error) {
	const op = "BookingService.server.grpcclient.RefundPayment"

	
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		c.logger.Error("Invalid create connection")
		return false, err
	}
	defer conn.Close()

	client := paymentV1.NewPaymentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.RefundPayment(ctx, &paymentV1.RefundPaymentRequest{
		PaymentId: paymendID,
	})
	if err != nil {
		c.logger.Error("Invalid grpc query",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)
		
		return false, err
	}

	return resp.Success, nil
}
