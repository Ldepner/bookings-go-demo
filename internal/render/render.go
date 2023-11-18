package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/ldepner/bookings/internal/config"
	"github.com/ldepner/bookings/internal/models"
)

var functions = template.FuncMap{}
var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer sets the config for the render package
func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

func Template(w http.ResponseWriter, r *http.Request, t string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	var err error
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			return err
		}
	}

	var tmpl *template.Template
	tmpl = tc[t]

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)
	if tmpl == nil {
		err = errors.New("could not find template")
		return err
	}

	err = tmpl.Execute(buf, td)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.html", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(fmt.Sprintf("%s/%s", pathToTemplates, name), fmt.Sprintf("%s/%s", pathToTemplates, "base.html"))

		if err != nil {
			return myCache, err
		}

		myCache[name] = ts
	}

	return myCache, nil
}
