package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/seed", app.seed)
	router.HandlerFunc(http.MethodGet, "/users", app.users)
	router.HandlerFunc(http.MethodPost, "/broadcast", app.broadcast)
	router.HandlerFunc(http.MethodPost, "/message", app.message)

	return router
}
