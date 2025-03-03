package main

import (
	"database/sql"
	"log"
)

func initSchema(db *sql.DB) error {
	createUser := `
		CREATE TABLE IF NOT EXISTS users (
		    id INTEGER PRIMARY KEY,
		    registration TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    user_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    bday DATE DEFAULT NULL,
		    username CHAR(64) UNIQUE NOT NULL,
		    full_name CHAR(255) UNIQUE NOT NULL,
		    email CHAR(255) UNIQUE NOT NULL,
		    phone CHAR(64) UNIQUE NOT NULL,
		    bio varchar(8192) UNIQUE NOT NULL,
		    passwd_hash CHAR(64) UNIQUE NOT NULL,
		    account_status INTEGER default 0
		);
	`

	_, err := db.Exec(createUser)
	if err != nil {
		log.Printf("initSchema() : create user table : %v", err)
		return err
	}

	createEmailV := `
		CREATE TABLE IF NOT EXISTS email_verification (
		    id INTEGER PRIMARY KEY,
		    userid INTEGER REFERENCES users(id),
		    verification_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    status INTEGER default 0,
		    token CHAR(64) UNIQUE NOT NULL
		);
	`

	_, err = db.Exec(createEmailV)
	if err != nil {
		log.Printf("initSchema() : create email verification table : %v", err)
		return err
	}

	createIncident := `
		CREATE TABLE IF NOT EXISTS incident (
		    id INTEGER PRIMARY KEY,
		    userid INTEGER REFERENCES users(id),
		    assigned_admin INTEGER REFERENCES users(id),
		    incident_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    status INTEGER default 0,
		    pending_action INTEGER DEFAULT 0,
		    incident_type INTEGER DEFAULT 0,
		    description VARCHAR(8192) UNIQUE NOT NULL
		);
	`

	_, err = db.Exec(createIncident)
	if err != nil {
		log.Printf("initSchema() : create incident table : %v", err)
		return err
	}

	return nil
}
