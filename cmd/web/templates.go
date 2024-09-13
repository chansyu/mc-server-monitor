package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"github.com/itzsBananas/mc-server-monitor/internal/mocks"
	"github.com/itzsBananas/mc-server-monitor/internal/models"
	"github.com/itzsBananas/mc-server-monitor/ui"
	"github.com/justinas/nosurf"
)

type templateData struct {
	Players         []string
	Response        models.Response
	Form            userLoginForm
	CSRFToken       string
	IsAuthenticated bool
	UserAutofill    struct{ Username, Password string }
	MockMode        bool
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("15:04")
}

func status(s bool) string {
	if s {
		return "Success"
	}
	return "Error"
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"status":    status,
}

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

func (app *application) newTemplateData(r *http.Request) *templateData {
	var userAutofill = struct{ Username, Password string }{}
	if app.mockMode {
		userAutofill.Username = mocks.Username
		userAutofill.Password = mocks.Password
	}
	return &templateData{
		CSRFToken:       nosurf.Token(r),
		IsAuthenticated: app.isAuthenticated(r),
		UserAutofill:    userAutofill,
		MockMode:        app.mockMode,
	}
}
