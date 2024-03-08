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

	players, err := app.rconConsole.Players()
	if err == nil {
		data.Players = players
	}

	app.renderPage(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) userLoginGet(w http.ResponseWriter, r *http.Request) {
	app.renderPage(w, http.StatusOK, "login.tmpl.html", nil)
}

func (app *application) seed(w http.ResponseWriter, r *http.Request) {
	seed, err := app.rconConsole.Seed()
	response := models.NewResponse("Seed", nil)

	if app.responseError(w, response, err) {
		return
	}
	app.responseSuccess(w, response, seed)
}

func (app *application) players(w http.ResponseWriter, r *http.Request) {
	players, err := app.rconConsole.Players()
	response := models.NewResponse("Players", nil)

	if app.responseError(w, response, err) {
		return
	}
	app.responseSuccess(w, response, strings.Join(players, ", "))
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

	if app.responseError(w, response, err) {
		return
	}
	app.responseSuccess(w, response, "Success!")
}

func (app *application) start(w http.ResponseWriter, r *http.Request) {
	err := app.adminConsole.Start(r.Context())
	response := models.NewResponse("Start", nil)

	if app.responseError(w, response, err) {
		return
	}
	app.responseSuccess(w, response, "Success!")
}

func (app *application) stop(w http.ResponseWriter, r *http.Request) {
	err := app.adminConsole.Stop(r.Context())
	response := models.NewResponse("Stop", nil)

	if app.responseError(w, response, err) {
		return
	}
	app.responseSuccess(w, response, "Success!")
}

func (app *application) restart(w http.ResponseWriter, r *http.Request) {
	err := app.adminConsole.Restart(r.Context())
	response := models.NewResponse("Restart", nil)

	if app.responseError(w, response, err) {
		return
	}
	app.responseSuccess(w, response, "Success!")
}

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	isReady, err := app.adminConsole.IsOnline(r.Context())

	response := models.NewResponse("Status", nil)

	if app.responseError(w, response, err) {
		return
	}

	msg := "Not Online!"
	if isReady {
		msg = "Online!"
	}
	app.responseSuccess(w, response, msg)
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
