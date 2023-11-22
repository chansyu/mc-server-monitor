package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.renderPage(w, http.StatusOK, "home.tmpl.html", nil)
}

func (app *application) seed(w http.ResponseWriter, r *http.Request) {
	output, err := app.remoteConsole.Seed()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Response: output}
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}

func (app *application) users(w http.ResponseWriter, r *http.Request) {
	output, err := app.remoteConsole.Users()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Response: strings.Join(output, "")}
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}

func (app *application) broadcast(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Message string `json:"message"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorLog.Fatalf(err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	output, err := app.remoteConsole.Broadcast(input.Message)
	if err != nil {
		app.serverError(w, err)
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
		app.clientError(w, http.StatusBadRequest)
		return
	}

	output, err := app.remoteConsole.Message(input.User, input.Message)
	if err != nil {
		app.serverError(w, err)
	} else {
		w.Write([]byte(output)) // No player was found or You whisper to itsBananas: this is a pm
	}
}
