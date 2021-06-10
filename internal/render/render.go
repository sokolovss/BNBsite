package render

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	config "github.com/sokolovss/BNBsite/internal/config"
	models "github.com/sokolovss/BNBsite/internal/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

var pathToTemplates = "./templates"

//NewTemplate gets app config for render package
func NewTemplate(a *config.AppConfig) {
	app = a
}

var functions = template.FuncMap{}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

//RenderTemplate is a template parser and executor
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, d *models.TemplateData) {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = NewTemplateCache()
		log.Println("UseCache = False. Rebuilding cache")

	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Error: no such key in template cache")
	}

	buf := new(bytes.Buffer)

	d = AddDefaultData(d, r)
	_ = t.Execute(buf, d)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error: sending respond to the client", err)
	}

}

//NewTemplateCache creates template cache as a map
func NewTemplateCache() (map[string]*template.Template, error) {
	pCache := make(map[string]*template.Template)
	p, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return pCache, err
	}

	for _, v := range p {
		n := filepath.Base(v)

		ts, err := template.New(n).Funcs(functions).ParseFiles(v)
		if err != nil {
			return pCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return pCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return pCache, err
			}
		}
		pCache[n] = ts
	}
	return pCache, nil
}
