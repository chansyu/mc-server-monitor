package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	console "github.com/itzsBananas/mc-server-monitor/internal/data"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	remoteConsole console.ConsoleInterface
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	rcaddr := flag.String("rcaddr", "127.0.0.1:25575", "Minecraft RCON API Address")
	rcpassword := flag.String("rcpass", "password", "Password to access Minecraft RCON API")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	con, err := console.Open(*rcaddr, *rcpassword)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		remoteConsole: con,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
