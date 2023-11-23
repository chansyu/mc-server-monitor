package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
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
		Message string `form:"message"`
	}

	err := app.decodePostForm(r, &input)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	output, err := app.remoteConsole.Broadcast(input.Message)
	if err != nil {
		app.serverError(w, err)
	} else {
		data := &templateData{Response: output}
		app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
	}
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

	output, err := app.remoteConsole.Message(input.User, input.Message)
	if err != nil {
		app.serverError(w, err)
	} else {
		data := &templateData{Response: output}
		app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
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
