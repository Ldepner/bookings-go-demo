package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ldepner/bookings/internal/driver"
	"github.com/ldepner/bookings/internal/helpers"
	"github.com/ldepner/bookings/internal/repository"
	"github.com/ldepner/bookings/internal/repository/dbrepo"
	"net/http"
	"strconv"
	"time"

	"github.com/ldepner/bookings/internal/forms"

	"github.com/ldepner/bookings/internal/config"
	"github.com/ldepner/bookings/internal/models"
	"github.com/ldepner/bookings/internal/render"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

var Repo *Repository

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (h *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	h.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Template(w, r, "home.html", &models.TemplateData{})
}

func (h *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "hello"

	remoteIP := h.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remoteIP"] = remoteIP

	render.Template(w, r, "about.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (h *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (h *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-02-01 -- 01/02 03:04:05PM `06 -0700

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "phone")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := h.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
		StartDate:     startDate,
		EndDate:       endDate,
	}

	err = h.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	h.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (h *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.html", &models.TemplateData{})
}

func (h *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (h *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (h *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.html", &models.TemplateData{})
}

func (h *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.html", &models.TemplateData{})
}

func (h *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.html", &models.TemplateData{})
}

func (h *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := h.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		h.App.ErrorLog.Print("Can't get error from session")
		h.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	h.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.html", &models.TemplateData{
		Data: data,
	})
}
