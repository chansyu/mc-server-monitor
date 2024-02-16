package admin_console

type AdminConsole interface {
	Start() error
	Restart() error
	Stop() error
	IsOnline() (bool, error)
}
