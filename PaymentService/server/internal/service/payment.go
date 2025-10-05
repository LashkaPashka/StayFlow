package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/config"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/lib/encode"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/model"
	"github.com/stripe/stripe-go/v78"
)

type PaymentClient interface {
	CreatePayment(payment *model.Payment) (*stripe.CheckoutSession, error)
	GetPaymentStatus(pID string) (*stripe.CheckoutSession, error)
}

type Storage interface {
	GetActivePayment(ctx context.Context, userID string) (payment model.Payment, err error)
	StoreCreatedPayment(ctx context.Context, payment *model.Payment) (err error)
	UpdatePayment(ctx context.Context, payment *model.Payment) (err error)
}

type Service struct {
	logger *slog.Logger
	config *config.Config
	PaymentClient
	Storage
}

func New(
	storage Storage,
	paymentClient PaymentClient,
	config *config.Config,
	logger *slog.Logger,
) *Service {
	return &Service{
		Storage: storage,
		PaymentClient: paymentClient,
		config: config,
		logger: logger,
	}
}

func (s *Service) CreatePayment(ctx context.Context, payment *model.Payment) (paymentID, url string,err error) {
	const op = "PaymentService.service.CreatePayemnt"
	
	publishKey := s.config.StripePublish
	
	// 1. Search user which exist payment
	existPayment, _ := s.Storage.GetActivePayment(ctx, payment.UserID)
	if len(existPayment.UserID) > 0 {
		s.logger.Info("payment exist",
			slog.String("pubkey", publishKey),
			slog.String("secret", existPayment.ClientSecret),
		)

		return "", "", errors.New("payment exist")
	}

	// 2. Create payment
	sessionResult, err := s.PaymentClient.CreatePayment(payment)
	if err != nil {
		s.logger.Error("error create payment using payemnt client",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return "", "", err
	}

	// 3. Update field clientSecret && paymentID in model Payment
	payment.ClientSecret = sessionResult.ClientSecret
	payment.PaymentID = sessionResult.ID
	
	// 4. Save model payment in Db
	if err = s.Storage.StoreCreatedPayment(ctx, payment); err != nil {
		s.logger.Error("Invalid store payment",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return "", "", err
	}
	
	return sessionResult.ID, sessionResult.URL, err
}

func (s *Service) VerifyPayment(ctx context.Context, payment *model.Payment) (err error) {
	const op = "PaymentService.service.VerifyPayment"

	existPayment, err := s.Storage.GetActivePayment(ctx, payment.UserID)
	if len(existPayment.UserID) > 0 || err != nil {
		s.logger.Error("no active payment exist")

		return errors.New("no exist payment")
	}

	paymentRes, err := s.PaymentClient.GetPaymentStatus(existPayment.PaymentID)
	if err != nil {
		s.logger.Error("error getPaymentStatus",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)
	}

	payment.Response = string(encode.Encode(paymentRes))
	payment.Status = model.PaymentStatusFailed

	if paymentRes.Status == "succeeded" {
		payment.Status = model.PaymentStatusConfirmed
	}

	if err = s.Storage.UpdatePayment(ctx, payment); err != nil {
		return err
	}

	return err
}