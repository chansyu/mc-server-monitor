package admin_console

import "context"

type AdminConsole interface {
	Start(context.Context) error
	Restart(context.Context) error
	Stop(context.Context) error
	IsOnline(context.Context) (bool, error)
	Close() error
}
