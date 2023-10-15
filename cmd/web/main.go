package main

import (
	"encoding/gob"
	"fmt"
	"github.com/ldepner/bookings/internal/models"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ldepner/bookings/internal/config"
	"github.com/ldepner/bookings/internal/handlers"
	"github.com/ldepner/bookings/internal/render"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("starting app on port %s", PORT))

	srv := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	// what am I going to store in my session?
	gob.Register(models.Reservation{})

	app.InProduction = false

	session = scs.New()
	app.Session = session
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	tc := map[string]*template.Template{}
	app.TemplateCache = tc
	app.UseCache = false
	render.NewTemplates(&app)

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	return nil
}
