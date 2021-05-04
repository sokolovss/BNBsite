package render

import (
	"bytes"
	"github.com/sokolovss/BNBsite/pkg/config"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

//NewTemplate gets app config for render package
func NewTemplate(a *config.AppConfig) {
	app = a
}

var functions = template.FuncMap{

}

//RenderTemplate is a template parser and executor
func RenderTemplate(w http.ResponseWriter, tmpl string) {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = NewTemplateCache()

	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Error: no such key in template cache")
	}

	buf := new(bytes.Buffer)

	_ = t.Execute(buf, nil)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error: sending respond to the client", err)
	}

}

//NewTemplateCache creates template cache as a map
func NewTemplateCache() (map[string]*template.Template, error) {
	pCache := make(map[string]*template.Template)
	p, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return pCache, err
	}

	for _, v := range p {
		n := filepath.Base(v)

		ts, err := template.New(n).Funcs(functions).ParseFiles(v)
		if err != nil {
			return pCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return pCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return pCache, err
			}
		}
		pCache[n] = ts
	}
	return pCache, nil
}
