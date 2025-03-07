package main

import "time"

type UserInstance struct {
	id             int
	registration   time.Time
	user_update    time.Time
	bday           time.Time
	username       string
	full_name      string
	email          string
	phone          string
	bio            string
	passwd_hash    string
	account_status int
}

type UserCredentials struct {
	id          int
	username    string
	email       string
	phone       string
	passwd_hash string
}

const (
	EmailUnverified = iota
	EmailVerified
)

const (
	AccountStatusActive = iota
	AccountStatusAdmin
	AccountStatusDisabled
	AccountStatusZombie
)

func AccountStatusString(status int) string {
	switch status {
	case AccountStatusActive:
		return "Active"
	case AccountStatusAdmin:
		return "Admin"
	case AccountStatusDisabled:
		return "Disabled"
	case AccountStatusZombie:
		return "Zombie"
	}
}
