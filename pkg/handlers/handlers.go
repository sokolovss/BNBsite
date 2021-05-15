package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sokolovss/BNBsite/pkg/config"
	"github.com/sokolovss/BNBsite/pkg/models"
	"github.com/sokolovss/BNBsite/pkg/render"
	"log"
	"net/http"
)

var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

//NewRepo creates new Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

//NewHandler set the Repository for the handlers
func NewHandler(r *Repository) {
	Repo = r
}

//Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, r, "index.page.tmpl", &models.TemplateData{})
}

//Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

//Colonels renders the room page
func (m *Repository) Colonels(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "colonels.page.tmpl", &models.TemplateData{})
}

//Contacts renders contacts page
func (m *Repository) Contacts(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contacts.page.tmpl", &models.TemplateData{})
}

//SearchAvailability renders search-availability page
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

//PostAvailability handles POST from search-availability form
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	s := r.Form.Get("start_date")
	e := r.Form.Get("end_date")
	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", s, e)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

//AvailabilityJSON handles request from availability form and returns JSON
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	resp := jsonResponse{
		OK: true,

		Message: "Available",
	}
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

//Reservation renders search-availability page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "reservation.page.tmpl", &models.TemplateData{})
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
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: sm,
	})

}
