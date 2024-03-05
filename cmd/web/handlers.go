package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.render(w, http.StatusOK, "home.tmpl.html", nil)
}

func (app *application) seed(w http.ResponseWriter, r *http.Request) {
	s, err := app.rconConsole.Seed()
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Write([]byte(s))
}

func (app *application) users(w http.ResponseWriter, r *http.Request) {
	s, err := app.rconConsole.Users()
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Write([]byte(fmt.Sprint(s)))
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

	output, err := app.rconConsole.Broadcast(input.Message)
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

	output, err := app.rconConsole.Message(input.User, input.Message)
	if err != nil {
		app.serverError(w, err)
	} else {
		w.Write([]byte(output)) // No player was found or You whisper to itsBananas: this is a pm
	}
}

func (app *application) start(w http.ResponseWriter, r *http.Request) {
	err := app.adminConsole.Start(r.Context())
	if err != nil {
		app.serverError(w, err)
	} else {
		w.Write([]byte("Success!"))
	}
}

func (app *application) stop(w http.ResponseWriter, r *http.Request) {
	err := app.adminConsole.Stop(r.Context())
	if err != nil {
		app.serverError(w, err)
	} else {
		w.Write([]byte("Success!"))
	}
}

func (app *application) restart(w http.ResponseWriter, r *http.Request) {
	err := app.adminConsole.Restart(r.Context())
	if err != nil {
		app.serverError(w, err)
	} else {
		w.Write([]byte("Success!"))
	}
}

func (app *application) ready(w http.ResponseWriter, r *http.Request) {
	isReady, err := app.adminConsole.IsOnline(r.Context())
	if err != nil {
		app.serverError(w, err)
	} else if !isReady {
		w.Write([]byte("Not Online!"))
	} else {
		w.Write([]byte("Online!"))
	}
}
