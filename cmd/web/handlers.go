package main

import (
	"html/template"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	ts, err := template.ParseFiles("./ui/html/pages/home.tmpl.html")
	if err != nil {
		app.errorLog.Print(err.Error())
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Print(err.Error())
		app.notFound(w)
	}
}
