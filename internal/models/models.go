package models

import "time"

// Reservation is the type for reservations in the system
type Reservation struct {
	ID		int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	Room      Room
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User is the type for users of the system
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	created_at  time.Time
	updated_at  time.Time
}

// Room is the type for rooms in the system
type Room struct {
	ID        int
	RoomName  string
	created_at time.Time
	updated_at time.Time
}

// Restriction is the type for restrictions in the system
type Restriction struct {
	ID        int
	RestrictionName  string
	created_at time.Time
	updated_at time.Time
}

// RoomRestriction is the type for room restrictions in the system
type RoomRestriction struct {
	ID            int
	RoomID        int
	RestrictionID int
	ReservationID int
	StartDate time.Time
	EndDate   time.Time
	Room	Room
	Reservations Reservation
	Restrictions Restriction

	created_at    time.Time
	updated_at    time.Time
}
