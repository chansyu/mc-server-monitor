package mocks

type NonAdminConsole struct{}

func (c NonAdminConsole) Seed() (string, error) {
	return "", nil
}
func (c NonAdminConsole) Broadcast(msg string) error {
	return nil
}
func (c NonAdminConsole) Message(user string, msg string) error {
	return nil
}
func (c NonAdminConsole) Players() ([]string, error) {
	return []string{}, nil
}