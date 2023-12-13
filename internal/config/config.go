package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/ldepner/bookings/internal/models"
	"html/template"
	"log"
)

// AppConfig holds the application config
type AppConfig struct {
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	UseCache      bool
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
