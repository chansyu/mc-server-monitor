package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/itzsBananas/mc-server-monitor/ui"
)

type templateData struct {
	Response string
}

var functions = template.FuncMap{}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	partials, err := fs.Glob(ui.Files, "html/partials/*.html")
	if err != nil {
		return nil, err
	}
	for _, partial := range partials {
		name := filepath.Base(partial)

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, partial)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}
