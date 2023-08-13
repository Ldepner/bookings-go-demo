package render

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ldepner/bookings/pkg/config"
	"github.com/ldepner/bookings/pkg/models"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, t string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if _, inMap := tc[t]; app.UseCache && !inMap {
		tc = app.TemplateCache
	} else {
		tc, _ = createTemplateCache(t)
	}

	var tmpl *template.Template
	tmpl = tc[t]

	td = AddDefaultData(td)
	_ = tmpl.Execute(w, td)
}

func createTemplateCache(t string) (map[string]*template.Template, error) {
	tc := app.TemplateCache
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.html",
	}

	tmpl, err := template.ParseFiles(templates...)

	if err != nil {
		return nil, err
	}

	tc[t] = tmpl
	return tc, nil
}
