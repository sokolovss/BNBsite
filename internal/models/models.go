package models

import "time"

//User is a model for users table
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//Room is a model for rooms table
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Restriction is a model for restrictions table
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//Reservation is a model for reservations table
type Reservation struct {
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
	Room      Room
}

//RoomRestriction is a model for room_restrictions table
type RoomRestriction struct {
	ID            int
	RoomID        int
	ReservationID int
	RestrictionID int
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}
