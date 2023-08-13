package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ldepner/bookings/pkg/config"
	"github.com/ldepner/bookings/pkg/handlers"
	"github.com/ldepner/bookings/pkg/render"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
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

	fmt.Println(fmt.Sprintf("starting app on port %s", PORT))

	srv := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}
