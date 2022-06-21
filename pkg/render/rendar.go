package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/varunkverma/bookings/pkg/config"
	"github.com/varunkverma/bookings/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders template using html/template
func RenderTemplate(w http.ResponseWriter, tmplName string, data *models.TemplateData) {

	var templateCache map[string]*template.Template
	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		var err error
		templateCache, err = CreateTemplateCache()
		if err != nil {
			log.Fatal("failed to create template cache")
		}
	}

	tmpl, ok := templateCache[tmplName]
	if !ok {
		log.Fatal("Couldn't get template from app config templateCache")
	}

	buf := new(bytes.Buffer)

	data = AddDefaultData(data)

	_ = tmpl.Execute(buf, data)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("error writing template to buffer", err.Error())
	}

}

// CreateTemplateCache creates a template map as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	log.Println("creating template cache...")
	// decalre the cache
	myCache := make(map[string]*template.Template)

	// find out paths to all the files that are in template folder ending with ".page.tmpl"
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return nil, err
	}
	// for each page, get its name and create a template with it
	for _, page := range pages {
		log.Println("processing template: ", page)
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		log.Println("looking our for layouts")
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return nil, err
		}
		log.Printf("%d layout(s) found\n", len(matches))

		// if there are any layouts present in the template folder, if the parsed template contains any it'll be parsed as well
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return nil, err
			}
		}

		myCache[name] = ts
	}
	log.Println("template cache created!")
	return myCache, nil
}
