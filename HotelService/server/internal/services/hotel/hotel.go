package hotel

import (
	"context"
	"log/slog"

	"github.com/LashkaPashka/HotelService/server/internal/config"
)

type Storage interface {
	GetAvailableRooms(ctx context.Context, hotelID, roomTypeID, checkIn, checkOut string) (availableRooms int, err error)
	BookedRoom(ctx context.Context, hotelID, roomTypeID, checkIn, checkOut string, rooms_count int) (success bool, err error)
	FreeRoom(ctx context.Context, hotelID, roomTypeID, checkIn, checkOut string, rooms_count int) (success bool, err error)
}

type HotelService struct {
	cfg     *config.Config
	logger  *slog.Logger
	storage Storage
}

func New(storage Storage, cfg *config.Config, logger *slog.Logger) *HotelService {
	return &HotelService{
		cfg:     cfg,
		logger:  logger,
		storage: storage,
	}
}

func (h *HotelService) CheckAvailability(
	ctx context.Context,
	hotelId,
	roomTypeID,
	checkIn,
	checkOut string,
) (available bool, availableRooms int, err error) {
	availableRooms, err = h.storage.GetAvailableRooms(ctx, hotelId, roomTypeID, checkIn, checkOut)
	if err != nil {
		h.logger.Error("Error storage of method GetAvailableRooms")
		return false, 0, err
	}

	if availableRooms == 0 {
		available = false
	}

	available = true

	return available, availableRooms, err
}

func (h *HotelService) ReserveRoom(
	ctx context.Context,
	hotelId,
	roomTypeID,
	checkIn,
	checkOut string,
	roomsCount int,
) (success bool, err error) {
	success, err = h.storage.BookedRoom(ctx, hotelId, roomTypeID, checkIn, checkOut, roomsCount)
	if err != nil {
		return false, err
	}

	return success, err
}

func (h *HotelService) ReleaseRoom(
	ctx context.Context,
	hotelID,
	roomTypeID,
	checkIn,
	checkOut string,
	roomsCount int,
) (success bool, err error) {
	success, err = h.storage.FreeRoom(ctx, hotelID, roomTypeID, checkIn, checkOut, roomsCount)
	if err != nil {
		return false, err
	}

	return success, err
}
