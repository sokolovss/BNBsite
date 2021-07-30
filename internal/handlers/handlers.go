package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	config "github.com/sokolovss/BNBsite/internal/config"
	"github.com/sokolovss/BNBsite/internal/driver"
	"github.com/sokolovss/BNBsite/internal/forms"
	"github.com/sokolovss/BNBsite/internal/helpers"
	models "github.com/sokolovss/BNBsite/internal/models"
	render "github.com/sokolovss/BNBsite/internal/render"
	"github.com/sokolovss/BNBsite/internal/repository"
	"github.com/sokolovss/BNBsite/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"time"
)

var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//NewRepo creates new Repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgreRepo(db.SQL, a),
	}
}

//NewHandler set the Repository for the handlers
func NewHandler(r *Repository) {
	Repo = r
}

//Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "index.page.tmpl", &models.TemplateData{})
}

//Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

//Colonels renders the room page
func (m *Repository) Colonels(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "colonels.page.tmpl", &models.TemplateData{})
}

//Contacts renders contacts page
func (m *Repository) Contacts(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contacts.page.tmpl", &models.TemplateData{})
}

//SearchAvailability renders search-availability page
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

//PostAvailability handles POST from search-availability form
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, r.Form.Get("start_date"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, r.Form.Get("end_date"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	availableRoom, err := m.DB.SearchAvailabilityAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(availableRoom) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = availableRoom

	res := models.Reservation{

		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}
	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/reservation", http.StatusSeeOther)
}

//BookRoom takes URL parameters. Builds sessional variable, redirects to make reservation page
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	//id, s , e parameters
	ID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	startDate := r.URL.Query().Get("s")
	endDate := r.URL.Query().Get("e")

	var res models.Reservation
	res.RoomID = ID
	layout := "2006-01-02"
	sd, _ := time.Parse(layout, startDate)
	ed, _ := time.Parse(layout, endDate)
	res.StartDate = sd
	res.EndDate = ed

	room, err := m.DB.SearchRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/reservation", http.StatusSeeOther)

}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

//AvailabilityJSON handles request from availability form and returns JSON
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	layout := "2006-01-02"

	sd := r.Form.Get("start")
	startDate, _ := time.Parse(layout, sd)

	ed := r.Form.Get("end")
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	fmt.Println(startDate, endDate, roomID)

	available, err := m.DB.SearchAvailabilityByDatesRoomID(startDate, endDate, roomID)
	if err != nil {
		helpers.ServerError(w, err)
	}

	resp := jsonResponse{
		OK:        available,
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
		Message:   "",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

//ReservationSummary handles reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation data")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

//Reservation renders search-availability page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	//var emptyReservation models.Reservation
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation out of the session"))
		return
	}
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	room, err := m.DB.SearchRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		StringMap: stringMap,
		Data:      data,
	})
}

//PostReservation handles posting of reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("can't get data from session"))
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3, r)
	form.MinLength("last_name", 2, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	reservID, err := m.DB.AddReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		RoomID:        reservation.RoomID,
		ReservationID: reservID,
		RestrictionID: 1,
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
	}

	err = m.DB.AddRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//perform logic
	sm := map[string]string{
		"test": "Hello again",
	}

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	sm["remote_ip"] = remoteIP

	//send to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: sm,
	})

}
