package console

import (
	"fmt"
	"strings"
	"time"

	"github.com/jltobler/go-rcon"
)

type ConsoleInterface interface {
	Users() ([]string, error)
	Seed() (string, error)
	Broadcast(msg string) (string, error)
	Message(user string, msg string) (string, error)
}

type Console struct {
	con     *rcon.Client
	timeout time.Duration
}

func Open(port string, password string, timeout time.Duration) *Console {
	return &Console{
		rcon.NewClient(port, password),
		timeout,
	}
}

func (c *Console) sendCommand(command string) (string, error) {
	success := make(chan string, 1)
	fail := make(chan error, 1)

	go func() {
		resp, err := c.con.Send(command)
		if err != nil {
			fail <- err
		} else {
			success <- resp
		}
	}()

	select {
	case resp := <-success:
		return resp, nil
	case err := <-fail:
		return "", err
	case <-time.After(c.timeout):
		return "", fmt.Errorf("console connection timeout")
	}
}

// "There are x users: \nBob, April\n" // what if no users
func (c *Console) Users() ([]string, error) {
	resp, err := c.sendCommand("list")
	if err != nil {
		return nil, err
	}

	str, err := stripPrefix(resp)
	if err != nil {
		return nil, err
	}
	str = strings.ReplaceAll(str, " ", "")
	list := strings.Split(str, ",")
	return list, nil
}

// "Seed: [1871644822592853811]"
func (c *Console) Seed() (string, error) {
	resp, err := c.sendCommand("seed")
	if err != nil {
		return "", err
	}

	if len(resp) < 9 {
		return "", fmt.Errorf("recieved malformed output from seed command: \"%s\"", resp)
	}

	return resp[7 : len(resp)-1], err
}

func (c *Console) Broadcast(msg string) (string, error) {
	command := fmt.Sprintf("/say %s", msg)
	resp, err := c.sendCommand(command)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (c *Console) Message(user string, msg string) (string, error) {
	command := fmt.Sprintf("/msg %s %s", user, msg)
	resp, err := c.sendCommand(command)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func stripPrefix(msg string) (string, error) {
	if idx := strings.IndexByte(msg, ':'); idx >= 0 {
		return msg[idx+1:], nil
	} else {
		return msg, fmt.Errorf("cannot strip prefix w/o colon: \"%s\"", msg)
	}
}
