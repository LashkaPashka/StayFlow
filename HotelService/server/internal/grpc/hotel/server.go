package hotelgrpc

import (
	"context"

	hotelsV1 "github.com/LashkaPashka/HotelService/hotel_proto/gen/go/hotel"
	"google.golang.org/grpc"
)

type Service interface {
	CheckAvailability(
		ctx context.Context,
		hotelId,
		roomTypeID,
		checkIn,
		checkOut string,
	) (available bool, availableRooms int, err error)

	ReserveRoom(
		ctx context.Context,
		hotelId,
		roomTypeID,
		checkIn,
		checkOut string,
		roomsCount int,
	) (success bool, err error)

	ReleaseRoom(
		ctx context.Context,
		hotelID,
		roomTypeID,
		checkIn,
		checkOut string,
		roomsCount int,
	) (success bool, err error)
}

type serverAPI struct {
	hotelsV1.UnimplementedHotelServiceServer
	hotelS Service
}

func Register(gRPCServer *grpc.Server, hotel Service) {
	hotelsV1.RegisterHotelServiceServer(gRPCServer, &serverAPI{hotelS: hotel})
}

func (s *serverAPI) CheckAvailability(
	ctx context.Context,
	in *hotelsV1.CheckAvailabilityRequest,
) (*hotelsV1.CheckAvailabilityResponse, error) {
	available, availableRooms, err := s.hotelS.CheckAvailability(ctx, in.HotelId, in.RoomTypeId, in.CheckIn, in.CheckOut)
	if err != nil {
		return &hotelsV1.CheckAvailabilityResponse{}, err
	}

	return &hotelsV1.CheckAvailabilityResponse{
		Available:      available,
		AvailableRooms: int32(availableRooms),
	}, nil
}

func (s *serverAPI) ReserveRoom(
	ctx context.Context,
	in *hotelsV1.ReserveRoomRequest,
) (*hotelsV1.ReserveRoomResponse, error) {
	success, err := s.hotelS.ReserveRoom(ctx, in.HotelId, in.RoomTypeId, in.CheckIn, in.CheckOut, int(in.RoomsCount))
	if err != nil {
		return &hotelsV1.ReserveRoomResponse{}, err
	}

	return &hotelsV1.ReserveRoomResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) ReleaseRoom(
	ctx context.Context,
	in *hotelsV1.ReleaseRoomRequest,
) (*hotelsV1.ReleaseRoomResponse, error) {
	success, err := s.hotelS.ReleaseRoom(ctx, in.HotelId, in.RoomTypeId, in.CheckIn, in.CheckOut, int(in.RoomsCount))
	if err != nil {
		return &hotelsV1.ReleaseRoomResponse{}, err
	}

	return &hotelsV1.ReleaseRoomResponse{
		Success: success,
	}, nil
}
