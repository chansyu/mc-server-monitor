package main

import (
	"net/http"

	"github.com/itzsBananas/mc-server-monitor/ui"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.FS(ui.Files))

	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/seed", app.seed)
	router.HandlerFunc(http.MethodGet, "/users", app.users)
	router.HandlerFunc(http.MethodPost, "/message", app.message)

	// container ctrl's (without authentication!!)
	// TODO: add authentication
	router.HandlerFunc(http.MethodPost, "/start", app.start)
	router.HandlerFunc(http.MethodPost, "/restart", app.restart)
	router.HandlerFunc(http.MethodPost, "/stop", app.stop)
	router.HandlerFunc(http.MethodPost, "/status", app.status)

	return router
}
