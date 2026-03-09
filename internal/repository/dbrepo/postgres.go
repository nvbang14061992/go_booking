package dbrepo

import (
	"context"
	"time"

	"github.com/bangn/bookings/internal/models"
)

func (m *PostgresDBRepo) AllUsers() bool {
	return true
}


// InsertReservation inserts a reservation into the database
func (m *PostgresDBRepo) InsertReservation(res models.Reservation) error {
	// create a context with timeout, to avoid long running queries, 
	// and potential memory leaks
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	_, err :=m.DB.ExecContext(
		ctx,
		stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}