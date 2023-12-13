package main

import (
	"encoding/gob"
	"fmt"
	"github.com/ldepner/bookings/internal/driver"
	"github.com/ldepner/bookings/internal/helpers"
	"github.com/ldepner/bookings/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ldepner/bookings/internal/config"
	"github.com/ldepner/bookings/internal/handlers"
	"github.com/ldepner/bookings/internal/render"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errLog *log.Logger

func main() {
	db, err := run()

	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)

	fmt.Println("starting mail listener...")
	ListenForMail()
	//
	//msg := models.MailData{
	//	To:      "blah@blah.com",
	//	From:    "me@here.com",
	//	Subject: "some subject",
	//	Content: "blah blah",
	//}
	//app.MailChan <- msg

	fmt.Println(fmt.Sprintf("starting app on port %s", PORT))

	srv := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to store in my session?
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to db
	log.Println("Connecting to database on port 5432...")
	db, err := driver.ConnectSql("host=localhost port=5432 dbname=bookings user=bytedance password=")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
		return nil, err
	}
	log.Println("Connected to database.")

	tc := map[string]*template.Template{}
	app.TemplateCache = tc
	app.UseCache = false
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	return db, nil
}
