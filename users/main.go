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
	Host   string `env:"HOST" envDefault:"localhost"`
	Port   string `env:"PORT" envDefault:"5432"`
	User   string `env:"USER" envDefault:"postgres"`
	Pass   string `env:"PASS"`
	DbName string `env:"DB_NAME" envDefault:"postgres"`
}

func initEnvConfig() (EnvConfig, error) {
	var c EnvConfig
	if c.Pass == "" {
		fmt.Println("EnvConfig::init() : missing password")
		return c, errors.New("missing password")
	}
	return c, nil
}

func (c *EnvConfig) getConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Pass, c.DbName)
}

func main() {
	// Generate a connection string
	conf, err := initEnvConfig()
	if err != nil {
		log.Fatal(err)
	}

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

	defer exit(db)
}

func exit(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("exit() : db close : %v", err)
	}
}
