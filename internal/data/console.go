package console

import (
	"fmt"
	"log"
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
	client *minecraft.Client
}

func Open(port string, password string) (*Console, error) {
	client, err := minecraft.NewClient(port)
	if err != nil {
		return nil, err
	}

	if err := client.Authenticate(password); err != nil {
		return nil, err
	}

	console := Console{
		client: client,
	}

	return &console, nil
}

func (c *Console) Users() ([]string, error) {
	resp, err := c.client.SendCommand("list")
	if err != nil {
		return nil, err
	}
	// log.Println(resp.Body) // "There are x users: \nBob, April\n" // what if no users

	str, err := stripPrefix(resp.Body)
	if err != nil {
		return nil, err
	}
	str = strings.ReplaceAll(str, " ", "")
	list := strings.Split(str, ",")
	return list, nil
}

func (c *Console) Seed() (string, error) {
	resp, err := c.client.SendCommand("seed")
	if err != nil {
		return "", err
	}
	log.Println(resp.Body) // "Seed: [1871644822592853811]"

	if len(resp.Body) < 9 {
		return "", fmt.Errorf("Recieved malformed output from seed command: \"%s\"", resp.Body)
	}

	return resp.Body[7 : len(resp.Body)-1], err
}

func (c *Console) Broadcast(msg string) (string, error) {
	command := fmt.Sprintf("/say %s", msg)
	resp, err := c.client.SendCommand(command)
	if err != nil {
		return "", err
	}
	return resp.Body, nil
}

func (c *Console) Message(user string, msg string) (string, error) {
	command := fmt.Sprintf("/msg %s %s", user, msg)
	resp, err := c.client.SendCommand(command)
	if err != nil {
		return "", err
	}
	return resp.Body, nil
}

func stripPrefix(msg string) (string, error) {
	if idx := strings.IndexByte(msg, ':'); idx >= 0 {
		return msg[idx+1:], nil
	} else {
		return msg, fmt.Errorf("Cannot strip prefix w/o colon: \"%s\"", msg)
	}
}
