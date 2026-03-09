package repository

import "github.com/bangn/bookings/internal/models"

//
type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (error)
}