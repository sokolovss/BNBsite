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
type User struct {
	ID           int
	FirstName    string
	LastName     string
	Email        string
	Password     string
	access_level int
	created_at   time.Time
	updated_at   time.Time
}

//Rooms is a model for users table
type Rooms struct {
	ID         int
	RoomName   string
	created_at time.Time
	updated_at time.Time
}
