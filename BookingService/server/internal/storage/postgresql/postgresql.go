package postgresql

import (
	"context"
	"log/slog"

	"github.com/LashkaPashka/BookingService/server/internal/config"
	"github.com/LashkaPashka/BookingService/server/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Pool   *pgxpool.Pool
	Logger *slog.Logger
}

func New(cfg config.Config, logger *slog.Logger) *Storage {
	const op = "bookingService.storage.postgresql.New"

	pool, err := pgxpool.New(context.Background(), cfg.StoragePath)
	if err != nil {
		logger.Error("Invalid connection Db", slog.String("File error: ", op))
		return nil
	}

	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("Invalid ping to Db", slog.String("File error:", op))
		return nil
	}

	logger.Info("Db was run successfuly!", slog.String("Storage_path", cfg.StoragePath))

	return &Storage{
		Pool:   pool,
		Logger: logger,
	}
}

func (s *Storage) SaveBooking(ctx context.Context, booking model.Booking) (booking_id string, err error) {
	const op = "bookingService.storage.postgresql.SaveBooking"

	query := `INSERT INTO bookings 
		(user_id, hotel_id, room_type_id, check_in, check_out, rooms_count, nights, total_amount, currency, status, idempotency_key)
		VALUES (@user_id, @hotel_id, @room_type_id, @check_in, @check_out, @rooms_count, @nights, @total_amount, @currency, @status, @idempotency_key)
		RETURNING booking_id`

	args := pgx.NamedArgs{
		"user_id":         booking.UserID,
		"hotel_id":        booking.HotelID,
		"room_type_id":    booking.RoomTypeID,
		"check_in":        booking.CheckIn,
		"check_out":       booking.CheckOut,
		"rooms_count":     booking.RoomsCount,
		"nights":          booking.Nights,
		"total_amount":    booking.TotalAmount,
		"currency":        booking.Currency,
		"status":          booking.Status,
		"idempotency_key": booking.IdempotencyKey,
	}

	err = s.Pool.QueryRow(ctx, query, args).Scan(&booking_id)
	if err != nil {
		s.Logger.Error("Invalid query SQL", slog.String("op", op))
		return "", err
	}

	s.Logger.Info("Booking was created successffuly in Db")

	return booking_id, err
}

func (s *Storage) UpdateStatusOfBooking(ctx context.Context, booking_id string, status string) (booking model.Booking, err error) {
	const op = "BookingService.storage.postgresql.UpdateStatusOfBooking"

	query := `UPDATE bookings 
			  SET status = @status 
			  WHERE booking_id = @booking_id
			  RETURNING user_id, hotel_id, room_type_id, TO_CHAR(check_in, 'YYYY-MM-DD'), 
			  			TO_CHAR(check_out, 'YYYY-MM-DD'), rooms_count, nights, status
			 `

	args := pgx.NamedArgs{
		"status":     status,
		"booking_id": booking_id,
	}

	err = s.Pool.QueryRow(ctx, query, args).Scan(
		&booking.UserID,
		&booking.HotelID,
		&booking.RoomTypeID,
		&booking.CheckIn,
		&booking.CheckOut,
		&booking.RoomsCount,
		&booking.Nights,
		&booking.Status,
	)
	if err != nil {
		s.Logger.Error("Error updated status", slog.String("op", op))
		return model.Booking{}, err
	}

	s.Logger.Info("Field status in bookingDb was updated successffuly")

	return booking, err
}

func (s *Storage) GetBooking(ctx context.Context, booking_id string, user_id string) (booking model.Booking, err error) {
	const op = "BookingService.storage.postgresql.GetBooking"

	query := `SELECT user_id, hotel_id, room_type_id, TO_CHAR(check_in, 'YYYY-MM-DD'), 
			  		 TO_CHAR(check_out, 'YYYY-MM-DD'), rooms_count, nights, status 
			  FROM bookings
			  WHERE (booking_id = @booking_id AND user_id = @user_id) 
			  ORDER BY nights
			  `

	args := pgx.NamedArgs{
		"booking_id": booking_id,
		"user_id":    user_id,
	}

	if err = s.Pool.QueryRow(ctx, query, args).Scan(
		&booking.UserID,
		&booking.HotelID,
		&booking.RoomTypeID,
		&booking.CheckIn,
		&booking.CheckOut,
		&booking.RoomsCount,
		&booking.Nights,
		&booking.Status,
	); err != nil {
		s.Logger.Error("Invalid getbooking query", slog.String("op", op))
		return model.Booking{}, err
	}

	return booking, err
}

func (s *Storage) GetBookings(ctx context.Context, user_id string) (bookings []model.Booking, err error) {
	const op = "BookingService.storage.postgresql.GetBookings"

	query := `SELECT user_id, hotel_id, room_type_id, TO_CHAR(check_in, 'YYYY-MM-DD'), 
			  		 TO_CHAR(check_out, 'YYYY-MM-DD'), rooms_count, nights, status 
			  FROM bookings
			  WHERE user_id = @user_id
			  ORDER BY nights
			  `
	args := pgx.NamedArgs{
		"user_id": user_id,
	}

	rows, err := s.Pool.Query(ctx, query, args)

	for rows.Next() {
		var booking model.Booking

		if err = rows.Scan(
			&booking.UserID,
			&booking.HotelID,
			&booking.RoomTypeID,
			&booking.CheckIn,
			&booking.CheckOut,
			&booking.RoomsCount,
			&booking.Nights,
			&booking.Status,
		); err != nil {
			s.Logger.Error("Invalid getbookings", slog.String("op", op))
			return bookings, err
		}

		bookings = append(bookings, booking)
	}

	return bookings, err
}
