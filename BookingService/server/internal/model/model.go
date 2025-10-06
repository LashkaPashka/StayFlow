package model

import "time"

type Booking struct {
	UserID string 			`json:"user_id"`
	HotelID string 			`json:"hotel_id"`
	RoomTypeID string 		`json:"room_type_id"`
	CheckIn string 			`json:"check_in"`
	CheckOut string 		`json:"check_out"`
	Nights int 				`json:"nights"`
	RoomsCount int 			`json:"rooms_count"`
	Status string 			`json:"status"`
	TotalAmount float64 	`json:"total_amount"`
	Currency string 		`json:"currency"`
	IdempotencyKey string 	`json:"idempotency_key"`
	CreatedAt time.Time 	`json:"created_at"`
	UpdatedAt time.Time 	`json:"updated_at"`
}

type PaymentInfo struct {
	TotalAmount		float64 	`json:"total_amount"`
	BookingID		string		`json:"booking_id"`
	UserID			string		`json:"user_id"`
	Currency 		string 		`json:"currency"`
	Method 			string 		`json:"method"`
	Token 			string 		`json:"token"`
	IdempotencyKey 	string 		`json:"idempotency_key"`
}