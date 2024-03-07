package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
	models "github.com/itzsBananas/mc-server-monitor/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	data := &templateData{}

	users, err := app.rconConsole.Users()
	if err == nil {
		data.Users = users
	}

	app.renderPage(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) seed(w http.ResponseWriter, r *http.Request) {
	seed, err := app.rconConsole.Seed()
	response := models.NewResponse("Seed", nil)

	if err != nil {
		app.serverError(w, err)
		response.ConsoleDisconnect()
	} else {
		response.ConsoleSuccess(seed)
	}

	data := &templateData{Response: *response}
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}

func (app *application) users(w http.ResponseWriter, r *http.Request) {
	users, err := app.rconConsole.Users()
	response := models.NewResponse("Users", nil)

	if err != nil {
		app.serverError(w, err)
		response.ConsoleDisconnect()
	} else {
		response.ConsoleSuccess(strings.Join(users, ", "))
	}

	data := &templateData{Response: *response}
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}

func (app *application) message(w http.ResponseWriter, r *http.Request) {
	var input struct {
		User    string `form:"user"`
		Message string `form:"message"`
	}

	err := app.decodePostForm(r, &input)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var response *models.Response
	if input.User == "All Players" {
		err = app.rconConsole.Broadcast(input.Message)
		response = models.NewResponse("Broadcast", []string{input.Message})
	} else {
		err = app.rconConsole.Message(input.User, input.Message)
		response = models.NewResponse("Message", []string{input.Message})
	}

	if err != nil {
		app.serverError(w, err)
		response.ConsoleDisconnect()
	} else {
		response.ConsoleSuccess("")
	}

	data := &templateData{Response: *response}
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)

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

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}
