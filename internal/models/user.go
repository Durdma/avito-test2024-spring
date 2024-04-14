package models

type User struct {
	Id      int  `json:"id"`
	TagId   int  `json:"tag_id"`
	IsAdmin bool `json:"is_admin"`
}
