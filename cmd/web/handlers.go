package main

import (
	"encoding/json"
	"fmt"
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

func (app *application) seed(w http.ResponseWriter, r *http.Request) {
	s, _ := app.remoteConsole.Seed()
	w.Write([]byte(s))
}

func (app *application) users(w http.ResponseWriter, r *http.Request) {
	s, _ := app.remoteConsole.Users()
	w.Write([]byte(fmt.Sprint(s)))
}

func (app *application) broadcast(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Message string `json:"message"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorLog.Fatalf(err.Error())
		return
	}

	output, err := app.remoteConsole.Broadcast(input.Message)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(output))
	}
}

func (app *application) message(w http.ResponseWriter, r *http.Request) {
	var input struct {
		User    string `json:"user"`
		Message string `json:"message"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorLog.Fatalf(err.Error())
		return
	}

	output, err := app.remoteConsole.Message(input.User, input.Message)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(output)) // No player was found or You whisper to itsBananas: this is a pm
	}
}
