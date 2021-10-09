package dbrepo

import (
	"context"
	"github.com/sokolovss/BNBsite/internal/models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

//AddReservation inserts reservation to database
func (m *postgresDBRepo) AddReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var resID int

	stmt := `insert into reservations (first_name,last_name,email,phone,start_date,end_date,room_id,
		created_at,updated_at) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&resID)
	if err != nil {
		return 0, err
	}

	return resID, nil
}

//AddRoomRestriction inserts a room restriction to DB
func (m *postgresDBRepo) AddRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id,reservation_id, restriction_id,
                               created_at, updated_at) 
		values ($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(
		ctx,
		stmt,
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

func (m *postgresDBRepo) SearchRoomByID(roomID int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	stmt := `select id,room_name,created_at,updated_at from rooms where id = $1`
	row := m.DB.QueryRowContext(
		ctx,
		stmt,
		roomID,
	)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil

}

//SearchAvailabilityByDatesRoomID returns TRUE if available
func (m *postgresDBRepo) SearchAvailabilityByDatesRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `select count(id) from room_restrictions where $1 < end_date and $2 > start_date and room_id = $3`

	row := m.DB.QueryRowContext(
		ctx,
		stmt,
		start,
		end,
		roomID,
	)
	var rowsNum int
	err := row.Scan(&rowsNum)

	if err != nil {
		return false, err
	}

	if rowsNum == 0 {
		return true, nil
	}
	return false, nil

}

// SearchAvailabilityAllRooms returns available rooms
func (m *postgresDBRepo) SearchAvailabilityAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	stmt := `select 
				r.id, r.room_name 
			from 
			     rooms r
			where r.id not in 
				(select room_id from room_restrictions where $1 < end_date and $2 > start_date)`

	rows, err := m.DB.QueryContext(
		ctx,
		stmt,
		start,
		end,
	)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

//GetUserByID returns user by id
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select id, first_name,last_name,email,password,access_level,created_at,updated_at from users where id = $1`
	rows := m.DB.QueryRowContext(
		ctx,
		query,
		id,
	)
	var u models.User
	err := rows.Scan(
		u.ID,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		u.CreatedAt,
		u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

//UpdateUser updates a user in te DB
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update users set first_name=$1, last_name=$2, email=$3,access_level=$4,updated_at=$5 where id = user`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}
