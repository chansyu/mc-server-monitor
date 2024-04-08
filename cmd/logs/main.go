package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hpcloud/tail"
)

func main() {
	logPath := getEnv("LOG_PATH", "./data/mc-server/logs/latest.log")

	t, err := tail.TailFile(logPath, tail.Config{Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}, Follow: true})
	if err != nil {
		log.Fatal(err)
	}
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
