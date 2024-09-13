package logs

import (
	"bufio"
	"fmt"
	"net"
)

type SocketInterface interface {
	AddClient(string) (<-chan string, error)
	RemoveClient(string) error
}

type Socket struct {
	addr    net.TCPAddr
	clients map[string]chan string
	conn    *net.TCPConn
	Logs    chan string
}

func OpenSocket(addr net.TCPAddr) *Socket {
	return &Socket{addr: addr, clients: make(map[string]chan string), Logs: make(chan string)}
}

func (s *Socket) AddClient(id string) (<-chan string, error) {
	if existingChannel, ok := s.clients[id]; ok {
		return existingChannel, fmt.Errorf("logsocket: client with %s already exist", id)
	}
	ch := make(chan string)
	s.clients[id] = ch
	if len(s.clients) == 1 {
		conn, err := s.connectLogs()
		if err != nil {
			return nil, err
		}
		s.conn = conn

	}
	return ch, nil
}

func (s *Socket) RemoveClient(id string) error {
	if _, ok := s.clients[id]; !ok {
		return fmt.Errorf("logsocket: client with %s doesn't exist", id)
	}
	close(s.clients[id])
	delete(s.clients, id)

	if len(s.clients) == 0 {
		err := s.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Socket) connectLogs() (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp", nil, &s.addr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to %v: %v", s.addr, err)
	}
	go func() {
		for connScanner := bufio.NewScanner(conn); connScanner.Scan(); {
			if err := connScanner.Err(); err != nil {
				s.Logs <- fmt.Sprintf("error reading from %s: %v", conn.RemoteAddr(), err)
				s.broadcast("Something went wrong with the logs...")
				conn.Close()
				return
			}

			s.broadcast(connScanner.Text())
		}
	}()
	return conn, nil
}

func (s *Socket) broadcast(msg string) {
	for _, cli := range s.clients {
		cli <- msg
	}
}
