package model

type Payment struct {
	UserID       string        `json:"user_id"`
	BookingID    string        `json:"booking_id"`
	Amount       float64       `json:"amount"`
	Currency     string        `json:"currency"`
	Method       string        `json:"method"`
	PaymentID    string        `json:"payment_id"`
	ClientSecret string        `json:"client_secret"`
	Status       PaymentStatus `json:"status"`
	Response     string        `json:"response"`
}

type PaymentStatus string

const (
	PaymentStatusInitial   PaymentStatus = "INITIAL"
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusConfirmed PaymentStatus = "CONFIRMED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
)
