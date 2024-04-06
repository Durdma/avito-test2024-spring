package models

import "time"

//TODO add auth manager

type Admin struct {
	Id           int
	RefreshToken string
	ExpiresAt    time.Time
}

type User struct {
	Id           int
	TagId        int
	RefreshToken string
	ExpiresAt    time.Time
}
