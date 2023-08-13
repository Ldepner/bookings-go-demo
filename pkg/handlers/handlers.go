package handlers

import (
	"net/http"

	"github.com/ldepner/bookings/pkg/config"
	"github.com/ldepner/bookings/pkg/models"
	"github.com/ldepner/bookings/pkg/render"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (h *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	h.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.html", &models.TemplateData{})
}

func (h *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "hello"

	remoteIP := h.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remoteIP"] = remoteIP

	render.RenderTemplate(w, "about.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
