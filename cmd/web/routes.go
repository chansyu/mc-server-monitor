package main

import (
	"net/http"

	"github.com/itzsBananas/mc-server-monitor/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.FS(ui.Files))

	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	basicMiddleware := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", basicMiddleware.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/seed", basicMiddleware.ThenFunc(app.seed))
	router.Handler(http.MethodGet, "/players", basicMiddleware.ThenFunc(app.players))
	router.Handler(http.MethodPost, "/message", basicMiddleware.ThenFunc(app.message))

	router.Handler(http.MethodGet, "/user/login", basicMiddleware.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", basicMiddleware.ThenFunc(app.userLoginPost))

	router.Handler(http.MethodGet, "/status", basicMiddleware.ThenFunc(app.status))

	protectedMiddleware := basicMiddleware.Append(app.requireAuthentication)

	router.Handler(http.MethodPost, "/user/logout", protectedMiddleware.ThenFunc(app.userLogoutPost))

	router.Handler(http.MethodPost, "/start", protectedMiddleware.ThenFunc(app.start))
	router.Handler(http.MethodPost, "/restart", protectedMiddleware.ThenFunc(app.restart))
	router.Handler(http.MethodPost, "/stop", protectedMiddleware.ThenFunc(app.stop))

	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
