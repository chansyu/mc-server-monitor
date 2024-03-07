package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/itzsBananas/mc-server-monitor/internal/console"
	"github.com/itzsBananas/mc-server-monitor/internal/models"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) renderPage(w http.ResponseWriter, status int, fileName string, data *templateData) {
	ts, ok := app.templateCache[fileName]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", fileName)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (app *application) renderPartial(w http.ResponseWriter, status int, fileName string, data *templateData) {
	ts, ok := app.templateCache[fileName]
	if !ok {
		err := fmt.Errorf("the partial template %s does not exist", fileName)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, baseName(fileName), data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func baseName(fileName string) string {
	f := strings.Split(fileName, ".")
	return f[0]
}

func (app *application) responseError(w http.ResponseWriter, response *models.Response, err error) bool {
	if err == nil {
		return false
	}

	if err != console.ErrTimeout {
		response.ConsoleDisconnect()
	} else {
		response.ConsoleError()
	}

	data := &templateData{Response: *response}
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)

	return true
}

func (app *application) responseSuccess(w http.ResponseWriter, response *models.Response, msg string) {
	response.ConsoleSuccess(msg)
	data := &templateData{Response: *response}
	app.renderPartial(w, http.StatusOK, "response.tmpl.html", data)
}
