package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/hpcloud/tail"
)

const PREFIX_FILTER = "[Server thread/INFO]:"

func main() {
	p := getEnv("LOG_PORT", "8081")
	port, err := strconv.Atoi(p)
	if err != nil {
		log.Fatal(err)
	}
	logPath := getEnv("LOG_PATH", "./data/mc-server/logs/latest.log")

	clients := NewClients()

	go func() {
		for {
			err = listenLogs(logPath, clients)
			if err != nil {
				log.Fatal(err)
			}
			<-time.After(2 * time.Second)
		}
	}()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Printf("listening at localhost: %s", listener.Addr())
	for {
		client, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go clients.handleClient(client)
	}
}

func listenLogs(logPath string, clients *Clients) error {
	t, err := tail.TailFile(logPath, tail.Config{Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}, Follow: true})
	if err != nil {
		return err
	}

	for line := range t.Lines {
		text := line.Text
		if len(text) < 35 {
			continue
		}
		if prefix := text[11:32]; prefix != PREFIX_FILTER {
			continue
		}
		msg := text[33:]
		log.Printf("Message sent: \"%s\"", msg)
		clients.broadcast(msg)
	}

	return nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

type Clients map[chan string]struct{}

func NewClients() *Clients {
	c := make(Clients)
	return &c
}

func (c *Clients) broadcast(data string) {
	for client := range *c {
		client <- data
	}
}

func (c *Clients) handleClient(client net.Conn) {
	log.Printf("%s - Client connected", client.RemoteAddr())
	eventChan := make(chan string)

	clientList := *c
	clientList[eventChan] = struct{}{}
	defer func() {
		delete(clientList, eventChan)
		close(eventChan)
		client.Close()
	}()

	for {
		data := <-eventChan
		_, err := fmt.Fprintf(client, "%s\n", data)
		if err != nil {
			log.Printf("%s - Client disconnected", client.RemoteAddr())
			return
		}
	}

}
