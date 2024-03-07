package console

import (
	"fmt"
	"strings"
	"time"

	models "github.com/itzsBananas/mc-server-monitor/internal/models"
	"github.com/jltobler/go-rcon"
)

type ConsoleInterface interface {
	Users() (*models.Response, error)
	Seed() (*models.Response, error)
	Broadcast(msg string) (*models.Response, error)
	Message(user string, msg string) (*models.Response, error)
}

type RCONConsole struct {
	con     *rcon.Client
	timeout time.Duration
}

func Open(port string, password string, timeout time.Duration) *RCONConsole {
	return &RCONConsole{
		rcon.NewClient(port, password),
		timeout,
	}
}

func (c *RCONConsole) sendCommand(command string) (string, error) {
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
func (c *RCONConsole) Users() (*models.Response, error) {
	reply, err := c.sendCommand("/list")
	resp := models.NewResponse("Users", nil)

	if err != nil {
		resp.ConsoleDisconnect()
		return resp, err
	}

	list, err := stripPrefix(reply)
	if err != nil {
		return resp, err
	}
	list = strings.ReplaceAll(list, " ", "")
	resp.ConsoleSuccess(list)
	return resp, nil
}

// "Seed: [1871644822592853811]"
func (c *RCONConsole) Seed() (*models.Response, error) {
	reply, err := c.sendCommand("/seed")
	resp := models.NewResponse("Seed", nil)

	if err != nil {
		resp.ConsoleDisconnect()
		return resp, err
	}

	if len(reply) < 9 {
		return resp, fmt.Errorf("recieved malformed output from seed command: \"%s\"", reply)
	}

	resp.ConsoleSuccess(reply[7 : len(reply)-1])
	return resp, err
}

func (c *RCONConsole) Broadcast(message string) (*models.Response, error) {
	command := fmt.Sprintf("/say %s", message)
	reply, err := c.sendCommand(command)
	resp := models.NewResponse("Broadcast Message", []string{message})

	if err != nil {
		resp.ConsoleDisconnect()
		return resp, err
	}
	resp.ConsoleSuccess(reply)
	return resp, nil
}

func (c *RCONConsole) Message(user string, message string) (*models.Response, error) {
	command := fmt.Sprintf("/msg %s %s", user, message)
	reply, err := c.sendCommand(command)
	resp := models.NewResponse("Private Message", []string{user, message})

	if err != nil {
		resp.ConsoleDisconnect()
		return resp, err
	}

	// TODO: need to check if user is not available
	if resp.Message == "No player was found" {
		return resp, err
	}
	resp.ConsoleSuccess(reply)
	return resp, nil
}

func stripPrefix(msg string) (string, error) {
	if idx := strings.IndexByte(msg, ':'); idx >= 0 {
		return msg[idx+1:], nil
	} else {
		return msg, fmt.Errorf("cannot strip prefix w/o colon: \"%s\"", msg)
	}
}
