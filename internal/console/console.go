package console

import (
	"fmt"
	"strings"

	"github.com/willroberts/minecraft-client"
)

type ConsoleInterface interface {
	Users() ([]string, error)
	Seed() (string, error)
	Broadcast(msg string) (string, error)
	Message(user string, msg string) (string, error)
}

type Console struct {
	port, password string
	client         *minecraft.Client
}

func Open(port string, password string) (*Console, error) {
	console := Console{
		port:     port,
		password: password,
	}

	client, err := console.newClient()
	console.client = client

	return &console, err
}

func (c *Console) newClient() (*minecraft.Client, error) {
	client, err := minecraft.NewClient(c.port)
	if err != nil {
		return nil, err
	}

	if err := client.Authenticate(c.password); err != nil {
		c.Close()
		return nil, err
	}

	return client, nil
}

func (c *Console) sendCommand(command string) (minecraft.Message, error) {
	if c.client == nil { // reconnect
		client, err := c.newClient()
		c.client = client
		if err != nil {
			return minecraft.Message{}, err
		}
	}
	resp, err := c.client.SendCommand(command)
	if err != nil {
		return minecraft.Message{}, err
	}
	return resp, nil
}

func (c *Console) Close() error {
	fmt.Println("closing")
	if c.client == nil {
		return nil
	}
	return c.Close()
}

// "There are x users: \nBob, April\n" // what if no users
func (c *Console) Users() ([]string, error) {
	resp, err := c.sendCommand("list")
	if err != nil {
		return nil, err
	}

	str, err := stripPrefix(resp.Body)
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

	if len(resp.Body) < 9 {
		return "", fmt.Errorf("recieved malformed output from seed command: \"%s\"", resp.Body)
	}

	return resp.Body[7 : len(resp.Body)-1], err
}

func (c *Console) Broadcast(msg string) (string, error) {
	command := fmt.Sprintf("/say %s", msg)
	resp, err := c.sendCommand(command)
	if err != nil {
		return "", err
	}
	return resp.Body, nil
}

func (c *Console) Message(user string, msg string) (string, error) {
	command := fmt.Sprintf("/msg %s %s", user, msg)
	resp, err := c.sendCommand(command)
	if err != nil {
		return "", err
	}
	return resp.Body, nil
}

func stripPrefix(msg string) (string, error) {
	if idx := strings.IndexByte(msg, ':'); idx >= 0 {
		return msg[idx+1:], nil
	} else {
		return msg, fmt.Errorf("cannot strip prefix w/o colon: \"%s\"", msg)
	}
}
