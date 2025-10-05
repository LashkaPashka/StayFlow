package payment

import (
	"log/slog"

	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/config"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/model"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

type Payment struct {
	stripeSecretKey  string
	successURL		 string
	cancelURL		 string
	logger 			 *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) *Payment {
	return &Payment{
		successURL: cfg.SuccessURL,
		cancelURL: cfg.CancelURL,
		stripeSecretKey: cfg.StripeSecret,
		logger: logger,
	}
}

func (p *Payment) CreatePayment(payment *model.Payment) (*stripe.CheckoutSession, error) {
	const op = "PaymentService.lib.stripe.CreatePayment"

	stripe.Key = p.stripeSecretKey

	amountInRub := payment.Amount * 100

	params := &stripe.CheckoutSessionParams{
	PaymentMethodTypes: stripe.StringSlice([]string{payment.Method}),
	LineItems: []*stripe.CheckoutSessionLineItemParams{
		{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				UnitAmount: stripe.Int64(int64(amountInRub)),
				Currency:   stripe.String(payment.Currency),

				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String("Hotel booking"),
					Description: stripe.String("Room reservation for your stay"),
				},
			},
			Quantity: stripe.Int64(1),
		},
	},
	Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
	SuccessURL: stripe.String(p.successURL),
	CancelURL:  stripe.String(p.cancelURL),
}



	params.AddMetadata("booking_id", payment.BookingID)
	params.AddMetadata("user_id", payment.UserID)

	session, err := session.New(params)
	if err != nil{
		p.logger.Error("Error creating payment session",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	return session, nil
}

func (p *Payment) GetPaymentStatus(pID string) (*stripe.CheckoutSession, error) {
	const op = "PaymentService.lib.stripe.GetPaymentStatus"
	
	stripe.Key = p.stripeSecretKey

	result, err := session.Get(pID, nil)

	if err != nil {
		p.logger.Error("Error getting payment session",
				slog.String("op", op),
				slog.String("err", err.Error()),
		)
		return nil, err
	}

	return result, err
}