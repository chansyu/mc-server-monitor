package serverStarter

// TODO: rename to have two implementations of a start/restart
type ClientInterface interface {
	Start() error
	Stop() error
	Restart() error
	Ready() bool
	Close() error
}
