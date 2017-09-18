package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	config Config
	db     *sql.DB
)

// Config Configuration data structure
type Config struct {
	DataDir string `json:"datadir"`

	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func main() {
	log.Println("Starting application...")

	data, err := ioutil.ReadFile("/etc/config.json")

	// Load configuration file
	switch {
	case os.IsNotExist(err):
		log.Println("Config file missing; using defaults")
		config = Config{
			DataDir:  "/var/lib/data",
			Host:     "127.0.0.1",
			Port:     "3306",
			Database: "test",
		}
	case err == nil:
		if err := json.Unmarshal(data, &config); err != nil {
			log.Fatal(err)
		}
	default:
		log.Println(err)
	}

	// Get additional config vars from environment
	log.Println("Overriding configuration from env vars.")
	if os.Getenv("APP_DATADIR") != "" {
		config.DataDir = os.Getenv("APP_DATADIR")
	}
	if os.Getenv("APP_HOST") != "" {
		config.DataDir = os.Getenv("APP_HOST")
	}
	if os.Getenv("APP_PORT") != "" {
		config.DataDir = os.Getenv("APP_PORT")
	}
	if os.Getenv("APP_USERNAME") != "" {
		config.DataDir = os.Getenv("APP_USERNAME")
	}
	if os.Getenv("APP_PASSWORD") != "" {
		config.DataDir = os.Getenv("APP_PASSWORD")
	}
	if os.Getenv("APP_DATABASE") != "" {
		config.DataDir = os.Getenv("APP_DATABASE")
	}

	// Load data directory, and create if it doesn't exist
	_, err = os.Stat(config.DataDir)
	if os.IsNotExist(err) {
		log.Println("Creating missing data directory", config.DataDir)
		err = os.MkdirAll(config.DataDir, 0755)
	}
	if err != nil {
		log.Fatal(err)
	}

	hostPort := net.JoinHostPort(config.Host, config.Port)
	log.Println("Connecting to database at", hostPort)

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=30s", config.Username, config.Password, hostPort, config.Database)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
	}

	// Try to connect to database until timeout is reached
	var dbError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		dbError = db.Ping()
		if dbError == nil {
			break
		}
		log.Println(dbError)
		time.Sleep(time.Duration(attempts) * time.Second)
	}

	if dbError != nil {
		log.Fatal(err)
	}
}
