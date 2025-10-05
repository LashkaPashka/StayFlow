package converter

import (
	paymentV1 "github.com/LashkaPashka/StayFlow/PaymentService/payment_proto/gen/go/payment"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/model"
)
 

func Convert(payload *paymentV1.CreatePaymentRequest) (payment *model.Payment) {
	return &model.Payment{
		UserID: payload.UserId,
		BookingID: payload.BookingId,
		Amount: payload.Amount,
		Currency: payload.Currency,
		Method: payload.Method,
		Status: model.PaymentStatusPending,

	}
}