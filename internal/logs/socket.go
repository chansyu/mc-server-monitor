package logs

import (
	"bufio"
	"fmt"
	"log"
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
}

func OpenSocket(addr net.TCPAddr) *Socket {
	return &Socket{addr: addr, clients: make(map[string]chan string)}
}

func (s *Socket) AddClient(id string) (<-chan string, error) {
	if existingChannel, ok := s.clients[id]; ok {
		return existingChannel, fmt.Errorf("logsocket: client with %s already exist", id)
	}
	ch := make(chan string)
	s.clients[id] = ch
	if len(s.clients) == 1 {
		conn, err := net.DialTCP("tcp", nil, &s.addr)
		if err != nil {
			log.Fatalf("error connecting to %v: %v", s.addr, err)
		}
		s.conn = conn
		go func() {
			for connScanner := bufio.NewScanner(conn); connScanner.Scan(); {
				for _, cli := range s.clients {
					cli <- connScanner.Text()
				}

				if err := connScanner.Err(); err != nil {
					log.Fatalf("error reading from %s: %v", conn.RemoteAddr(), err)
				}
				if connScanner.Err() != nil {
					log.Fatalf("error reading from %s: %v", conn.RemoteAddr(), err)
				}
			}
		}()
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
