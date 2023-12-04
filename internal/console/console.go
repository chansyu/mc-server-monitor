package console

import (
	"fmt"
	"strings"
	"time"

	"github.com/jltobler/go-rcon"
)

type ConsoleInterface interface {
	Users() (*Response, error)
	Seed() (*Response, error)
	Broadcast(msg string) (*Response, error)
	Message(user string, msg string) (*Response, error)
}

type ConsoleModel struct {
	con     *rcon.Client
	timeout time.Duration
}

func Open(port string, password string, timeout time.Duration) *ConsoleModel {
	return &ConsoleModel{
		rcon.NewClient(port, password),
		timeout,
	}
}

func (c *ConsoleModel) sendCommand(command string) (string, error) {
	success := make(chan string, 1)
	fail := make(chan error, 1)

	go func() {
		reply, err := c.con.Send(command)
		if err != nil {
			fail <- err
		} else {
			success <- reply
		}
	}()

	select {
	case reply := <-success:
		return reply, nil
	case err := <-fail:
		return "", err
	case <-time.After(c.timeout):
		return "", fmt.Errorf("console connection timeout")
	}
}

// "There are x users: \nBob, April\n" // what if no users
func (c *ConsoleModel) Users() (*Response, error) {
	reply, err := c.sendCommand("/list")
	resp := newResponse("Users", nil)

	if err != nil {
		resp.consoleDisconnect()
		return resp, err
	}

	list, err := stripPrefix(reply)
	if err != nil {
		return resp, err
	}
	list = strings.ReplaceAll(list, " ", "")
	resp.consoleSuccess(list)
	return resp, nil
}

// "Seed: [1871644822592853811]"
func (c *ConsoleModel) Seed() (*Response, error) {
	reply, err := c.sendCommand("/seed")
	resp := newResponse("Seed", nil)

	if err != nil {
		resp.consoleDisconnect()
		return resp, err
	}

	if len(reply) < 9 {
		return resp, fmt.Errorf("recieved malformed output from seed command: \"%s\"", reply)
	}

	resp.consoleSuccess(reply[7 : len(reply)-1])
	return resp, err
}

func (c *ConsoleModel) Broadcast(message string) (*Response, error) {
	command := fmt.Sprintf("/say %s", message)
	reply, err := c.sendCommand(command)
	resp := newResponse("Broadcast Message", []string{message})

	if err != nil {
		resp.consoleDisconnect()
		return resp, err
	}
	resp.consoleSuccess(reply)
	return resp, nil
}

func (c *ConsoleModel) Message(user string, message string) (*Response, error) {
	command := fmt.Sprintf("/msg %s %s", user, message)
	reply, err := c.sendCommand(command)
	resp := newResponse("Private Message", []string{user, message})

	if err != nil {
		resp.consoleDisconnect()
		return resp, err
	}

	// TODO: need to check if user is not available
	if resp.Message == "No player was found" {
		return resp, err
	}
	resp.consoleSuccess(reply)
	return resp, nil
}

func stripPrefix(msg string) (string, error) {
	if idx := strings.IndexByte(msg, ':'); idx >= 0 {
		return msg[idx+1:], nil
	} else {
		return msg, fmt.Errorf("cannot strip prefix w/o colon: \"%s\"", msg)
	}
}
