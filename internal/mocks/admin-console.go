package mocks

import "context"

type AdminConsole struct{}

func (c AdminConsole) Start(context.Context) error {
	return nil
}
func (c AdminConsole) Restart(context.Context) error {
	return nil
}
func (c AdminConsole) Stop(context.Context) error {
	return nil
}
func (c AdminConsole) IsOnline(context.Context) (bool, error) {
	return true, nil
}
func (c AdminConsole) Close() error {
	return nil
}
