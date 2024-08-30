package console

import (
	"fmt"
	"strings"
	"time"

	"github.com/jltobler/go-rcon"
)

type NonAdmin interface {
	Players() ([]string, error)
	Seed() (string, error)
	Broadcast(msg string) error
	Message(user string, msg string) error
}

type RCON struct {
	con     *rcon.Client
	timeout time.Duration
}

func Open(port string, password string, timeout time.Duration) *RCON {
	return &RCON{
		rcon.NewClient(port, password),
		timeout,
	}
}

func (c *RCON) sendCommand(command string) (string, error) {
	success := make(chan string, 1)
	fail := make(chan error, 1)

	go func() {
		reply, err := c.con.Send(command)
		if err != nil {
			fail <- fmt.Errorf("%w: %v", ErrInternal, err)
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
		return "", ErrTimeout
	}
}

// "There are x users: Bob, April\n" // what if no users
func (c *RCON) Players() ([]string, error) {
	reply, err := c.sendCommand("/list")

	if err != nil {
		return nil, err
	}

	players, err := stripPrefix(reply)
	if err != nil {
		return nil, err
	}

	players = strings.Join(strings.Fields(players), "")
	if len(players) == 0 {
		return nil, nil
	}
	list := strings.Split(players, ",")
	return list, nil
}

// "Seed: [1871644822592853811]"
func (c *RCON) Seed() (string, error) {
	reply, err := c.sendCommand("/seed")

	if err != nil {
		return "", err
	}

	if len(reply) < 9 {
		return "", ErrMalformedOutput
	}

	return reply[7 : len(reply)-1], nil
}

func (c *RCON) Broadcast(message string) error {
	command := fmt.Sprintf("/say %s", message)
	_, err := c.sendCommand(command)
	return err
}

func (c *RCON) Message(user string, message string) error {
	command := fmt.Sprintf("/msg %s %s", user, message)
	reply, err := c.sendCommand(command)

	if err != nil {
		return err
	}

	// TODO: need to check if user is not available
	if reply == "No player was found" {
		return ErrNoPlayer
	}
	return nil
}

func stripPrefix(msg string) (string, error) {
	if idx := strings.IndexByte(msg, ':'); idx >= 0 {
		return msg[idx+1:], nil
	} else {
		return msg, ErrMalformedOutput
	}
}
