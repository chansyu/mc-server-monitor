package admin_console

type AdminConsoleInterface interface {
	Start() error
	Restart() error
	Stop() error
	IsOnline() (bool, error)
	Close()
}
