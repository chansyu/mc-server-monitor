package console

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorcon/rcon"
)

type NonAdmin interface {
	Players() ([]string, error)
	Seed() (string, error)
	Broadcast(msg string) error
	Message(user string, msg string) error
}

type RCON struct {
	port     string
	password string
	timeout  time.Duration
}

func Open(port string, password string, timeout time.Duration) *RCON {
	return &RCON{
		port, password, timeout,
	}
}

func (c *RCON) sendCommand(command string) (string, error) {
	conn, err := rcon.Dial(c.port, c.password, rcon.SetDeadline(c.timeout))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	resp, err := conn.Execute(command)
	if err != nil {
		return "", err
	}

	return resp, nil
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
