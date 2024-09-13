package mocks

import (
	"fmt"
	"time"
)

type LogsSocket struct {
	clients map[string]client
}

type client struct {
	cli  chan string
	done chan string
}

func OpenLogsSocket() LogsSocket {
	return LogsSocket{make(map[string]client)}
}

func (s *LogsSocket) AddClient(id string) (<-chan string, error) {
	if existingClient, ok := s.clients[id]; ok {
		return existingClient.cli, fmt.Errorf("logsocket: client with %s already exist", id)
	}
	c := client{make(chan string), make(chan string)}

	go func() {
		for {
			select {
			case <-time.After(2 * time.Second):
				c.cli <- "hey there"
			case <-c.done:
				return
			}
		}
	}()
	s.clients[id] = c
	return c.cli, nil
}

func (s *LogsSocket) RemoveClient(id string) error {
	if _, ok := s.clients[id]; !ok {
		return fmt.Errorf("logsocket: client with %s doesn't exist", id)
	}
	close(s.clients[id].cli)
	close(s.clients[id].done)
	delete(s.clients, id)
	return nil
}
