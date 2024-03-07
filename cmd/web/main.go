package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/form/v4"
	admin_console "github.com/itzsBananas/mc-server-monitor/internal/admin-console"
	console "github.com/itzsBananas/mc-server-monitor/internal/console"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
	rconConsole   console.ConsoleInterface
	adminConsole  admin_console.AdminConsole
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

	rcon := console.Open(rconAddress, rconPassword, rconTimeout)

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	adminConsole, err := getAdminConsole()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer adminConsole.Close()

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		rconConsole:   rcon,
		templateCache: templateCache,
		adminConsole:  adminConsole,
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

func getAdminConsole() (admin_console.AdminConsole, error) {
	mode := getEnv("MODE", "production")

	var adminConsole admin_console.AdminConsole
	var err error
	if mode == "production" {
		gcpProject := getEnv("GCP_PROJECT", "PROJECT_NAME")
		gcpZone := getEnv("GCP_ZONE", "ZONE_NAME")
		gcpInstance := getEnv("GCP_INSTANCE", "INSTANCE_NAME")

		adminConsole, err = admin_console.GCPAdminConsoleOpen(gcpProject, gcpInstance, gcpZone)
	} else {
		localContainerId := getEnv("LOCAL_CONTAINER_ID", "mc-server")
		adminConsole, err = admin_console.LocalAdminConsoleOpen(localContainerId)
	}

	if err != nil {
		return nil, err
	}
	return adminConsole, nil
}
