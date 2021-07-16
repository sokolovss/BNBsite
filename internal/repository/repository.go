package repository

import (
	"github.com/sokolovss/BNBsite/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool

	AddReservation(res models.Reservation) (int, error)

	AddRoomRestriction(r models.RoomRestriction) error

	SearchAvailabilityByDatesRoomID(start, end time.Time, roomID int) (bool, error)

	SearchAvailabilityAllRooms(start, end time.Time) ([]models.Room, error)

	SearchRoomByID(roomID int) (models.Room, error)
}
