package postgresql

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
	logger *slog.Logger
}

func New(connStr string, logger *slog.Logger) *Storage {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		logger.Error("Invalid connection to Db")
		return nil
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Error("Error ping to Db")
		return nil
	}

	logger.Info("Connect Db was made successfully!")

	return &Storage{
		pool: pool,
		logger: logger,
	}
}

func (s *Storage) GetAvailableRooms(
	ctx context.Context,
	hotelID, 
	roomTypeID, 
	checkIn, 
	checkOut string,
) (availableRooms int, err error){
	const op = "HotelService.storage.postgresql.GetAvailableRooms"

	query := `SELECT MIN(available) 
			  FROM room_inventory
			  WHERE (hotel_id = @hotel_id
				  AND room_type_id = @room_type_id 
				  AND date BETWEEN @check_in AND @check_out)
			 `

	args := pgx.NamedArgs{
		"hotel_id": hotelID,
		"room_type_id": roomTypeID,
		"check_in": checkIn,
		"check_out": checkOut,
	}
	
	err = s.pool.QueryRow(ctx, query, args).Scan(&availableRooms)
	if err != nil {
		s.logger.Error("Invalid sql query",
			slog.String("err", err.Error()),
			slog.String("op", op),
		)

		return availableRooms, err
	}

	s.logger.Info("GetAvailableRooms was completed successfully!")

	return availableRooms, err
}

func (s *Storage) BookedRoom(
		ctx context.Context,
		hotelID, 
		roomTypeID, 
		checkIn, 
		checkOut string,
		roomsCount int,
	) (success bool, err error) {
	const op = "HotelService.storage.postgresql.BookedRoom"

	query := `UPDATE room_inventory
			  SET available = available - @rooms_count,
				  reserved = reserved + @rooms_count
			  WHERE (hotel_id = @hotel_id 
				     AND room_type_id = @room_type_id
					 AND date BETWEEN @check_in AND @check_out)
			`
	args := pgx.NamedArgs{
		"hotel_id": hotelID,
		"room_type_id": roomTypeID,
		"check_in": checkIn,
		"check_out": checkOut,
		"rooms_count": roomsCount,
	}

	if _, err = s.pool.Exec(ctx, query, args); err != nil {
		s.logger.Error("Invalid sql query",
					    slog.String("err", err.Error()),
						slog.String("op", op),
		)
		return false, err
	}

	s.logger.Info("BookedRoom was completed successfully!")

	return true, err
}

func (s *Storage) FreeRoom(
		ctx context.Context, 
		hotelID, 
		roomTypeID, 
		checkIn, 
		checkOut string, 
		roomsCount int,
	) (success bool, err error) {
	const op = "HotelService.storage.postgresql.FreeRoom"

	query := `UPDATE room_inventory
			  SET available = available + @rooms_count,
				  reserved = reserved - @rooms_count
			  WHERE (hotel_id = @hotel_id 
				     AND room_type_id = @room_type_id
					 AND date BETWEEN @check_in AND @check_out)
			 `
	args := pgx.NamedArgs{
		"hotel_id": hotelID,
		"room_type_id": roomTypeID,
		"check_in": checkIn,
		"check_out": checkOut,
		"rooms_count": roomsCount,
	}

	if _, err = s.pool.Exec(ctx, query, args); err != nil {
		s.logger.Error("Invalid sql query",
					    slog.String("err", err.Error()),
						slog.String("op", op),
		)
		return false, err
	}

	s.logger.Info("FreeRoom was completed successfully!")

	return true, err
}