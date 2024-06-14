package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	admin_console "github.com/itzsBananas/mc-server-monitor/internal/admin-console"
	console "github.com/itzsBananas/mc-server-monitor/internal/console"
	"github.com/itzsBananas/mc-server-monitor/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	rconConsole    console.ConsoleInterface
	adminConsole   admin_console.AdminConsole
	sessionManager *scs.SessionManager
	users          *models.UserModel
	mcLogs         *LogsSocket
}

func main() {
	serverAddress := getEnv("SERVER_ADDRESS", ":8080")
	rconAddress := getEnv("RCON_ADDRESS", "rcon://127.0.0.1:25575")
	rconPassword := getEnv("RCON_PASSWORD", "password")
	rconTimeoutString := getEnv("RCON_TIMEOUT", "5s")
	dsn := getEnv("DSN", "file:./data/mc-server-monitor.db?_timeout=5000")
	logsAddress := getEnv("LOGS_ADDRESS", "127.0.0.1:8081")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	ip, err := net.ResolveTCPAddr("tcp", logsAddress)
	if err != nil {
		errorLog.Fatal(err)
	}
	logsSocket := OpenLogsSocket(*ip)

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

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatalf(dsn)
	}
	defer db.Close()

	sessionManager := scs.New()
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	adminConsole, err := getAdminConsole()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer adminConsole.Close()

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		rconConsole:    rcon,
		templateCache:  templateCache,
		adminConsole:   adminConsole,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		users:          &models.UserModel{DB: db},
		mcLogs:         logsSocket,
	}

	// Removed timeouts to support SSE seamlessly
	srv := &http.Server{
		Addr:     serverAddress,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		// IdleTimeout: time.Minute,
		// ReadTimeout:  5 * time.Second,
		// WriteTimeout: 10 * time.Second,
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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

type LogsSocket struct {
	addr    net.TCPAddr
	clients map[string]chan string
	conn    *net.TCPConn
}

func OpenLogsSocket(addr net.TCPAddr) *LogsSocket {
	return &LogsSocket{addr: addr, clients: make(map[string]chan string)}
}

func (s *LogsSocket) AddClient(id string) (<-chan string, error) {
	if existingChannel, ok := s.clients[id]; ok {
		return existingChannel, fmt.Errorf("logsocket: client with %s already exist", id)
	}
	ch := make(chan string)
	s.clients[id] = ch
	if len(s.clients) == 1 {
		conn, err := net.DialTCP("tcp", nil, &s.addr)
		if err != nil {
			log.Fatalf("error connecting to %v: %v", s.addr, err)
		}
		s.conn = conn
		go func() {
			for connScanner := bufio.NewScanner(conn); connScanner.Scan(); {
				for _, cli := range s.clients {
					cli <- connScanner.Text()
				}

				if err := connScanner.Err(); err != nil {
					log.Fatalf("error reading from %s: %v", conn.RemoteAddr(), err)
				}
				if connScanner.Err() != nil {
					log.Fatalf("error reading from %s: %v", conn.RemoteAddr(), err)
				}
			}
		}()
	}
	return ch, nil
}

func (s *LogsSocket) RemoveClient(id string) error {
	if _, ok := s.clients[id]; !ok {
		return fmt.Errorf("logsocket: client with %s doesn't exist", id)
	}
	close(s.clients[id])
	delete(s.clients, id)

	if len(s.clients) == 0 {
		err := s.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
