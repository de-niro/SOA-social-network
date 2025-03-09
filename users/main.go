package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Config interface {
	getConnString() string
}

type EnvConfig struct {
	ServerHost string `env:"SERVER_HOST" envDefault:"localhost"`
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER" envDefault:"postgres"`
	DBPass     string `env:"DB_PASS"`
	DBName     string `env:"DB_NAME" envDefault:"postgres"`
}

func initEnvConfig() (EnvConfig, error) {
	var c EnvConfig
	if c.DBPass == "" {
		fmt.Println("EnvConfig::init() : missing password, variable $DB_PASS is unset")
		return c, errors.New("missing password")
	}
	return c, nil
}

func (c *EnvConfig) getConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName)
}

func main() {
	conf, err := initEnvConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Generate a connection string
	connString := conf.getConnString()
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("main() : db conn : %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("main() : db ping : %v", err)
	}

	err = initSchema(db)
	if err != nil {
		log.Fatalf("main() : init db schema : %v", err)
	}

	hostAddr := concat(conf.ServerHost, ":", conf.ServerPort)
	log.Printf("main() : serving at %v", hostAddr)

	execLoopIgnition(db, hostAddr)
	defer exit(db)
}

func exit(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("exit() : db close : %v", err)
	}
}
