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
func (m *PostgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// create a context with timeout, to avoid long running queries, 
	// and potential memory leaks
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int
	
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err :=m.DB.QueryRowContext(
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
	).Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *PostgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) (error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}


// SearchAvailabilityByDatesByRoomId
func (m *PostgresDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var numRows int

	query := `
	SELECT
		COUNT(id)
	FROM
		room_restrictions
	WHERE
		room_id = $1 AND
		$2 < end_date and $3 > start_date;`

	row := m.DB.QueryRowContext(
		ctx,
		query,
		roomId,
		start,
		end,
	)
	err := row.Scan(&numRows)
	if err != nil {
		return false, nil
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}


// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *PostgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room

	query := `
	SELECT
		r.id, r.room_name
	FROM
		roomS r
	WHERE r.id not in
		(SELECT room_id FROM room_restrictions rr WHERE $1 < rr.end_date AND $2 > rr.start_date);`

	rows, err := m.DB.QueryContext(
		ctx,
		query,
		start,
		end,
	)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	return  rooms, nil
}

// GetRoomByID gets a room by ID
func (m *PostgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var room models.Room

	query := `
		SELECT 
			id, room_name, created_at, updated_at
		FROM
			rooms
		WHERE
			id = $1`
	
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}