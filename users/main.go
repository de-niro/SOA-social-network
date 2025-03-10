package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Netflix/go-env"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type Config interface {
	getConnString() string
}

type EnvConfig struct {
	ServerHost    string `env:"SERVER_HOST,default=localhost"`
	ServerPort    string `env:"SERVER_PORT,default=8080"`
	DBHost        string `env:"DB_HOST,default=localhost"`
	DBPort        string `env:"DB_PORT,default=5432"`
	DBUser        string `env:"DB_USER,default=postgres"`
	DBPass        string `env:"DB_PASS"`
	DBName        string `env:"DB_NAME,default=postgres"`
	DBConnRetries int    `env:"DB_CONN_RETRIES,default=5"`
}

func getPrintableEnvConfig(config EnvConfig) EnvConfig { config.DBPass = ""; return config }

func initEnvConfig() (EnvConfig, error) {
	var c EnvConfig
	_, err := env.UnmarshalFromEnviron(&c)
	if err != nil {
		log.Fatal(err)
	}

	if c.DBPass == "" {
		fmt.Println("EnvConfig::init() : missing password, variable $DB_PASS is unset")
		return c, errors.New("missing password")
	}

	fmt.Printf("initEnvConfig() : initializing server with parameters %+v\n", getPrintableEnvConfig(c))

	return c, nil
}

func (c *EnvConfig) getConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName)
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

	success := false
	for i := 0; i < conf.DBConnRetries; i++ {
		err = db.Ping()
		if err != nil {
			time.Sleep(time.Second)
		} else {
			success = true
		}
	}
	if !success {
		log.Fatalf("main() : db retries exceeded : %v", err)
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
