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
	response := models.NewResponse("Start", nil)

	if err != nil {
		app.serverError(w, err)
		response.ConsoleDisconnect()
	} else {
		response.ConsoleSuccess("")
	}

	data := &templateData{Response: *response}
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}

func (app *application) stop(w http.ResponseWriter, r *http.Request) {
	err := app.adminConsole.Stop(r.Context())
	response := models.NewResponse("Stop", nil)

	if err != nil {
		app.serverError(w, err)
		response.ConsoleDisconnect()
	} else {
		response.ConsoleSuccess("")
	}

	data := &templateData{Response: *response}
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}

func (app *application) restart(w http.ResponseWriter, r *http.Request) {
	err := app.adminConsole.Restart(r.Context())
	response := models.NewResponse("Restart", nil)

	if err != nil {
		app.serverError(w, err)
		response.ConsoleDisconnect()
	} else {
		response.ConsoleSuccess("")
	}

	data := &templateData{Response: *response}
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}

func (app *application) isOnline(w http.ResponseWriter, r *http.Request) {
	isReady, err := app.adminConsole.IsOnline(r.Context())

	response := models.NewResponse("Status", nil)

	if err != nil {
		app.serverError(w, err)
		response.ConsoleDisconnect()
	} else {
		if isReady {
			response.ConsoleSuccess("Online!")
		} else {
			response.ConsoleSuccess("Not Online!")
		}
	}

	data := &templateData{Response: *response}
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
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
