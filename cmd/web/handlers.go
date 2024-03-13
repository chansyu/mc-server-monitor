package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
	models "github.com/itzsBananas/mc-server-monitor/internal/models"
)

const MsgSuccessDefault = "Succeeded!"

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

	if err != nil {
		app.consoleError(err, response)
	} else {
		response.Succeed(seed)
	}

	data := app.newTemplateData(r)
	data.Response = response
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}

func (app *application) players(w http.ResponseWriter, r *http.Request) {
	players, err := app.rconConsole.Players()
	response := models.NewResponse("Players", nil)

	if err != nil {
		app.consoleError(err, response)
	} else {
		response.Succeed(strings.Join(players, ", "))
	}

	data := app.newTemplateData(r)
	data.Response = response
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

	var response models.Response
	if input.User == "All Players" {
		err = app.rconConsole.Broadcast(input.Message)
		response = models.NewResponse("Broadcast", []string{input.Message})
	} else {
		err = app.rconConsole.Message(input.User, input.Message)
		response = models.NewResponse("Message", []string{input.Message})
	}

	if err != nil {
		app.consoleError(err, response)
	} else {
		response.Succeed(MsgSuccessDefault)
	}

	data := app.newTemplateData(r)
	data.Response = response
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}

func (app *application) start(w http.ResponseWriter, r *http.Request) {
	err := app.adminConsole.Start(r.Context())
	response := models.NewResponse("Start", nil)

	if err != nil {
		app.consoleError(err, response)
	} else {
		response.Succeed(MsgSuccessDefault)
	}

	data := app.newTemplateData(r)
	data.Response = response
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}

func (app *application) stop(w http.ResponseWriter, r *http.Request) {
	err := app.adminConsole.Stop(r.Context())
	response := models.NewResponse("Stop", nil)

	if err != nil {
		app.consoleError(err, response)
	} else {
		response.Succeed("Succeeded")
	}

	data := app.newTemplateData(r)
	data.Response = response
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}

func (app *application) restart(w http.ResponseWriter, r *http.Request) {
	err := app.adminConsole.Restart(r.Context())
	response := models.NewResponse("Restart", nil)

	if err != nil {
		app.consoleError(err, response)
	} else {
		response.Succeed(MsgSuccessDefault)
	}

	data := app.newTemplateData(r)
	data.Response = response
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	isReady, err := app.adminConsole.IsOnline(r.Context())

	response := models.NewResponse("Status", nil)

	if err != nil {
		app.consoleError(err, response)
	} else {
		msg := "Not Online!"
		if isReady {
			msg = "Online!"
		}
		response.Succeed(msg)
	}

	data := app.newTemplateData(r)
	data.Response = response
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
