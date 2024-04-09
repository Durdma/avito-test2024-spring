package models

import "time"

//TODO add auth manager

type User struct {
	Id           int
	TagId        int
	IsAdmin      bool
	RefreshToken string
	ExpiresAt    time.Time
}
