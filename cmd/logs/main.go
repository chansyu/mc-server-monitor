package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/hpcloud/tail"
)

func main() {
	p := getEnv("LOG_PORT", "8081")
	port, err := strconv.Atoi(p)
	if err != nil {
		log.Fatal(err)
	}
	logPath := getEnv("LOG_PATH", "./data/mc-server/logs/latest.log")

	t, err := tail.TailFile(logPath, tail.Config{Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}, Follow: true})
	if err != nil {
		log.Fatal(err)
	}

	clients := NewClients()

	go func() {
		for line := range t.Lines {
			clients.broadcast(line.Text)
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
	println("Client connected")
	eventChan := make(chan string)

	clientList := *c
	clientList[eventChan] = struct{}{}
	defer func() {
		delete(clientList, eventChan)
		close(eventChan)
	}()

	for {
		data := <-eventChan
		println("Sending data to client", data)
		_, err := fmt.Fprintf(client, "data: %s\n", data)
		if err != nil {
			fmt.Println("Client disconnected")
			return
		}
	}

}
