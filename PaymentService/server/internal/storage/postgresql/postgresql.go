package postgresql

import (
	"context"
	"log/slog"

	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
	looger *slog.Logger
}

func New(connStr string, logger *slog.Logger) *Storage {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		logger.Error("Invalid connection to Db",
			   slog.String("err", err.Error()),
		)
		return nil
	}

	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("Invalid ping to Db",
				slog.String("err", err.Error()),
		)
		return nil
	}

	return &Storage{
		pool: pool,
		looger: logger,
	}
}

func (s *Storage) StoreCreatedPayment(ctx context.Context, payment *model.Payment) (err error) {
	const op = "PaymentService.Storage.CreatePayment"
	
	query := `INSERT INTO payments
			  (user_id, booking_id, amount, currency, method, payment_id, client_secret, status, response)
			  VALUES (@user_id, @booking_id, @amount, @currency, @method, @payment_id, @client_secret, @status, @response) 
			`

	args := pgx.NamedArgs{
		"user_id": 			payment.UserID,
		"booking_id": 		payment.BookingID,
		"currency":			payment.Currency,
		"method":			payment.Method,
		"payment_id":		payment.PaymentID,
		"amount": 			payment.Amount,
		"client_secret": 	payment.ClientSecret,
		"status": 			payment.Status,
		"response":			payment.Response,
	}

	n, err := s.pool.Exec(ctx, query, args)
	if err != nil {
		s.looger.Error("Invalid createPayment",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)
		return err
	}

	s.looger.Info("Data was added successfully in Payment", slog.Int64("quantity lines", n.RowsAffected()))

	return err
}

func (s *Storage) UpdatePayment(ctx context.Context, payment *model.Payment) (err error) {
	const op = "PaymentService.Storage.UpdatePayment"
	
	query := `UPDATE payments
			  SET status = @status,
				  response = @response
			  WHERE (user_id = @user_id
			  		AND booking_id = @booking_id)
			`

	args := pgx.NamedArgs{
		"status": payment.Status,
		"response": payment.Response,
		"user_id": payment.UserID,
		"booking_id": payment.BookingID,
	}

	n, err := s.pool.Exec(ctx, query, args)
	if err != nil {
		s.looger.Error("Invalid UpdatePayment",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return err
	}

	s.looger.Info("Fields status && response was updated succesffuly in Db", slog.Int64("Update field", n.RowsAffected()))

	return err
}

func (s *Storage) GetActivePayment(ctx context.Context, userID string) (payment model.Payment, err error) {
	const op = "PaymentService.Storage.GetActivePayment"

	
	query := `SELECT user_id, booking_id, amount, currency, method, payment_id, client_secret, status, response
			  FROM payments
			  WHERE user_id = @user_id AND status = @status
			  ORDER BY created_at DESC
			  LIMIT 1
	`

	args := pgx.NamedArgs{
		"user_id": userID,
		"status":  model.PaymentStatusPending,
	}

	err = s.pool.QueryRow(ctx, query, args).Scan(
		&payment.UserID,
		&payment.BookingID,
		&payment.Amount,
		&payment.Currency,
		&payment.Method,
		&payment.PaymentID,
		&payment.ClientSecret,
		&payment.Status,
		&payment.Response,
	)
	if err != nil {
		s.looger.Error("Invalid getActivePayment",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return payment, err
	}
	
	s.looger.Info("Payment was given successfully from Db")

	return payment, err
}