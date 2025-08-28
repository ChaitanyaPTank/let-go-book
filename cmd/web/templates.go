package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.chaitanya.observer/internal/models"
	"snippetbox.chaitanya.observer/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("2 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// init map
	cache := map[string]*template.Template{}

	// Glob returns slice of all matching files according to pattern
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		// parse from the embedded fs
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// add template set to the map with name of the page as the key
		cache[name] = ts
	}
	return cache, nil
}
