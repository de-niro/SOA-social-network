package main

import (
	"database/sql"
	"errors"
	"log"
	"time"
	api "users/users_api"
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

func getUserIDByUsername(db *sql.DB, username string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&id)
	return id, err
}

func getUserIDByCredentials(db *sql.DB, credentials UserCredentials) (int, error) {
	var id int
	if credentials.username != "" {
		err := db.QueryRow("SELECT id FROM users WHERE username=$1", credentials.username).Scan(&id)
		return id, err
	} else if credentials.email != "" {
		err := db.QueryRow("SELECT id FROM users WHERE email=$1", credentials.email).Scan(&id)
		return id, err
	} else if credentials.phone != "" {
		err := db.QueryRow("SELECT id FROM users WHERE phone=$1", credentials.phone).Scan(&id)
		return id, err
	} else {
		return -1, errors.New("missing credentials")
	}
}

func getUserInfoByUsername(db *sql.DB, username string) (UserInstance, error) {
	var user UserInstance
	var passwd string
	err := db.QueryRow("SELECT * FROM users WHERE username=$1", username).Scan(&user.id, &user.registration, &user.user_update, &user.bday, &user.username, &user.full_name, &user.email, &user.phone, &user.bio, &user.passwd_hash, &passwd, &user.account_status)
	return user, err
}

func getUserCredentials(db *sql.DB, id int) (UserCredentials, error) {
	var user UserCredentials
	err := db.QueryRow("SELECT id, username, email, phone, passwd_hash FROM users WHERE id=$1", id).Scan(&user.id, &user.username, &user.email, &user.phone, &user.passwd_hash)
	return user, err
}

func isUserEmailVerified(db *sql.DB, id int) (bool, error) {
	var status int
	rows, err := db.Query("SELECT status FROM email_verification ev JOIN users u ON ev.userid = u.id WHERE u.id=$1", id)
	defer rows.Close()
	if err != nil {
		return false, err
	}

	for rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			return false, err
		}
		if status == EmailVerified {
			return true, nil
		}
	}
	err = rows.Err()
	return false, err
}
