package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

// AppConfig holds the application config
type AppConfig struct {
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	UseCache      bool
	InProduction  bool
	Session       *scs.SessionManager
}
