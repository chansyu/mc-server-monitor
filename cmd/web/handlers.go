package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/form/v4"
	models "github.com/itzsBananas/mc-server-monitor/internal/models"
	validator "github.com/itzsBananas/mc-server-monitor/internal/validator"
)

const MsgSuccessDefault = "Succeeded!"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	data := app.newTemplateData(r)

	players, err := app.rconConsole.Players()
	if err == nil {
		data.Players = players
	}

	app.renderPage(w, http.StatusOK, "home.tmpl.html", data)
}

type userLoginForm struct {
	Username            string `form:"username"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.renderPage(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Username), "Username", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "Password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.infoLog.Println(data.Form)
		app.renderPage(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Username, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.renderPage(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

func (app *application) logsGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.renderPage(w, http.StatusOK, "logs.tmpl.html", data)
}

func (app *application) logsSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	rc := http.NewResponseController(w)

	for i := 0; i < 10; i++ {
		fmt.Fprintf(w, "data: <p>%s</p>\n\n", fmt.Sprintf("Event %d", i))
		time.Sleep(2 * time.Second)
		rc.SetWriteDeadline(time.Now().Add(5 * time.Second))
		err := rc.Flush()
		if err != nil {
			app.infoLog.Print(err)
		}
	}
}
