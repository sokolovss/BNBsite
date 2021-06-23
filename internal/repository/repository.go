package repository

import "github.com/sokolovss/BNBsite/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	AddReservation(res models.Reservation) error
}
