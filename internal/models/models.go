package models

import "time"

//Reservation model for reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

//Users is a model for users table
type Users struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//Rooms is a model for users table
type Rooms struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Restrictions is a model for restrictions table
type Restrictions struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//Reservations is a model for reservations table
type Reservations struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	RoomID    int
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

//RoomRestrictions is a model for room_restrictions table
type RoomRestrictions struct {
	ID            int
	RoomID        int
	ReservationID int
	RestrictionID int
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
