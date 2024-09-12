package mocks

import "fmt"

type LogsSocket struct {
	clients map[string]chan string
}

func OpenLogsSocket() LogsSocket {
	return LogsSocket{make(map[string]chan string)}
}

func (s *LogsSocket) AddClient(id string) (<-chan string, error) {
	if existingChannel, ok := s.clients[id]; ok {
		return existingChannel, fmt.Errorf("logsocket: client with %s already exist", id)
	}
	ch := make(chan string)
	s.clients[id] = ch

	return ch, nil
}

func (s *LogsSocket) RemoveClient(id string) error {
	if _, ok := s.clients[id]; !ok {
		return fmt.Errorf("logsocket: client with %s doesn't exist", id)
	}
	close(s.clients[id])
	delete(s.clients, id)
	return nil
}
