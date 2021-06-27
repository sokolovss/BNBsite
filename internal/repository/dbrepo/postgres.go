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

func (m * postgresDBRepo) SearchAvailabilityByDates(start, end time.Time) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt:= `select count(id) form room_restrictions where $1 < end_date and $2 > start_date `



	row := m.DB.QueryRowContext(
		ctx,
		stmt,
		start,
		end,
	)
	var rowsNum int
	err := row.Scan(&rowsNum)

	if err != nil {
		return rowsNum,err
	}


	return rowsNum, nil
}