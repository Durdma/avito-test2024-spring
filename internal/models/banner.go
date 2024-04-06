package models

import "time"

type Banner struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`

	Tags    []Tag   `json:"tags"`
	Feature Feature `json:"feature"`

	IsActive bool `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Feature struct {
	ID int `json:"feature_id"`
}

type Tag struct {
	ID int `json:"tags_id"`
}
