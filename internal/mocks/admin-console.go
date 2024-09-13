package mocks

import "context"

type AdminConsole struct{ isOnline bool }

func (c *AdminConsole) Start(context.Context) error {
	c.isOnline = true
	return nil
}
func (c *AdminConsole) Restart(context.Context) error {
	c.isOnline = true
	return nil
}
func (c *AdminConsole) Stop(context.Context) error {
	c.isOnline = false
	return nil
}
func (c *AdminConsole) IsOnline(context.Context) (bool, error) {
	return c.isOnline, nil
}
func (c *AdminConsole) Close() error {
	return nil
}
