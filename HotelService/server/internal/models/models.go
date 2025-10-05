package models

import "time"

type Hotel struct {
	CreatedAt	time.Time	`json:"created_at"`
	Name		string		`json:"name"`
	Address		string		`json:"address"`
	City		string		`json:"city"`
	Country		string		`json:"country"`
}

type RoomTypes struct {
	BasePrice	float64		`json:"base_price"`
	Capacity	int			`json:"capacity"`
	HotelID 	string		`json:"hotel_id"`
	Code		string		`json:"code"`
}

type RoomInventory struct {
	Date		time.Time	`json:"date"`
	Available	int			`json:"available"`
	Reserved	int			`json:"reserved"`
	HotelID		string		`json:"hotel_id"`
	RoomTypeID	string		`json:"room_type_id"`
}