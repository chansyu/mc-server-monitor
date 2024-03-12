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

	basicMiddleWare := func(handler func(http.ResponseWriter, *http.Request)) http.Handler {
		return app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(handler)))
	}

	router.Handler(http.MethodGet, "/", basicMiddleWare(app.home))
	router.Handler(http.MethodGet, "/seed", basicMiddleWare(app.seed))
	router.Handler(http.MethodGet, "/users", basicMiddleWare(app.players))
	router.Handler(http.MethodPost, "/message", basicMiddleWare(app.message))

	router.HandlerFunc(http.MethodGet, "/user/login", app.userLoginGet)

	// container ctrl's (without authentication!!)
	// TODO: add authentication
	router.Handler(http.MethodPost, "/start", basicMiddleWare(app.start))
	router.Handler(http.MethodPost, "/restart", basicMiddleWare(app.restart))
	router.Handler(http.MethodPost, "/stop", basicMiddleWare(app.stop))
	router.Handler(http.MethodGet, "/status", basicMiddleWare(app.status))

	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
