package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

var regexURL = regexp.MustCompile(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)

type Banner struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

type AdminBanner struct {
	ID      int     `json:"id"`
	Content Banner  `json:"content"`
	Tags    []Tag   `json:"tags_ids"`
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

func (b *Banner) ValidateBanner() error {
	if len(b.Title) == 0 {
		return errors.New(fmt.Sprintf("title is empty"))
	}

	if len(b.Title) > 255 {
		return errors.New(fmt.Sprintf("title length is too long."+
			" it must be less than %v, but have %v", 255, len(b.Title)))
	}

	if len(b.Text) == 0 {
		return errors.New(fmt.Sprintf("text is empty"))
	}

	if len(b.Text) > 1000 {
		return errors.New(fmt.Sprintf("text length is too long."+
			" it must be less than %v, but have %v", 255, len(b.Text)))
	}

	if len(b.URL) == 0 {
		return errors.New(fmt.Sprintf("url is empty"))
	}

	if !regexURL.MatchString(b.URL) {
		return errors.New("url is incorrect")
	}

	return nil
}

func (b *AdminBanner) ValidateAndSetTags(tags []int) error {
	for _, t := range tags {
		if t < 0 {
			return errors.New(fmt.Sprintf("tag id must be greater or equal to 0, but have %v", t))
		}

		b.Tags = append(b.Tags, Tag{ID: t})
	}

	return nil
}

func (b *AdminBanner) ValidateAndSetFeature(feature int) error {
	if feature < 0 {
		return errors.New(fmt.Sprintf("feature id must be greater or equal to 0, but have %v", feature))
	}

	b.Feature.ID = feature

	return nil
}
