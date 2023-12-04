package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
	console "github.com/itzsBananas/mc-server-monitor/internal/console"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	output, _ := app.console.Users()
	var users []string
	if len(output.Message) > 0 {
		users = strings.Split(output.Message, ",")
	}

	data := &templateData{Users: users}
	app.renderPage(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) seed(w http.ResponseWriter, r *http.Request) {
	output, err := app.console.Seed()
	if err != nil {
		app.serverError(w, err)
	}

	data := &templateData{Response: *output}
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}

func (app *application) users(w http.ResponseWriter, r *http.Request) {
	output, err := app.console.Users()
	if err != nil {
		app.serverError(w, err)
	}

	data := &templateData{Response: *output}
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

	var output *console.Response
	if input.User == "All Players" {
		output, err = app.console.Broadcast(input.Message)
	} else {
		output, err = app.console.Message(input.User, input.Message)
	}

	if err != nil {
		app.serverError(w, err)
	}
	data := &templateData{Response: *output}
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
