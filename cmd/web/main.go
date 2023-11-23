package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/form/v4"
	console "github.com/itzsBananas/mc-server-monitor/internal/console"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	remoteConsole console.ConsoleInterface
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	serverAddress := getEnv("SERVER_ADDRESS", ":8080")
	rconAddress := getEnv("RCON_ADDRESS", "rcon://127.0.0.1:25575")
	rconPassword := getEnv("RCON_PASSWORD", "password")
	rconTimeoutString := getEnv("RCON_TIMEOUT", "5s")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	rconTimeout, err := time.ParseDuration(rconTimeoutString)
	if err != nil {
		rconTimeout = 5 * time.Second
	}

	con := console.Open(rconAddress, rconPassword, rconTimeout)

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		remoteConsole: con,
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	srv := &http.Server{
		Addr:     serverAddress,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", serverAddress)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
