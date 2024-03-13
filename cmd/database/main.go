package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type myScanner struct {
	*bufio.Scanner
}

type myDB struct {
	*sql.DB
}

const MsgYes = "Y"

func main() {
	dsn := getEnv("DSN", "file:./mc-server-monitor.db?_timeout=5000")
	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	scanner := myScanner{bufio.NewScanner(os.Stdin)}
	fmt.Println("Database User Script:")
	for {
		input := scanner.mustReadInput("Create user ('Y' for yes, otherwise for no): ")
		if input != MsgYes {
			return
		}
		username := scanner.mustReadInput("Username: ")
		password := scanner.mustReadInput("Password: ")
		db.createUser(username, password)
	}

}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func openDB(dsn string) (myDB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return myDB{}, err
	}
	if err = db.Ping(); err != nil {
		return myDB{}, err
	}
	return myDB{db}, nil
}

func (scanner myScanner) mustReadInput(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return scanner.Text()
}

func (db myDB) execSqlScript(fileName string) {
	script, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		log.Fatal(err)
	}
}

func (db myDB) createUser(username string, password string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatal(err)
	}
	stmt := `INSERT INTO users (username, hashed_password, created)
VALUES(?, ?, CURRENT_TIMESTAMP)`
	_, err = db.Exec(stmt, username, string(hashedPassword))
	if err != nil {
		log.Fatal(err)
	}
}
