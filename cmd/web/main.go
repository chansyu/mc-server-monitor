package main

import (
	"log"
	"net/http"
	"os"

	"github.com/itzsBananas/mc-server-monitor/internal/config"
	console "github.com/itzsBananas/mc-server-monitor/internal/data"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	remoteConsole console.ConsoleInterface
}

func main() {
	conf := config.New()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	con, err := console.Open(conf.RconAddress, conf.RconPassword)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		remoteConsole: con,
	}

	srv := &http.Server{
		Addr:     conf.ServerAddress,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", conf.ServerAddress)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
