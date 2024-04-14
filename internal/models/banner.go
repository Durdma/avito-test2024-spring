package models

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"time"
)

var regexURL = regexp.MustCompile(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)

type Banner struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

type AdminBanner struct {
	ID      int     `json:"banner_id"`
	Content Banner  `json:"content"`
	Tags    []Tag   `json:"tags_ids"`
	Feature Feature `json:"feature_id"`

	IsActive bool `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Feature struct {
	ID int `json:"feature_id"`
}

type Tag struct {
	ID int `json:"tag_id"`
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
	if len(slices.Compact(tags)) != len(tags) {
		return errors.New("list of tags_ids contain similar ids")
	}

	for _, t := range tags {
		if t < 0 {
			return errors.New(fmt.Sprintf("tag id must be greater or equal to 0, but have %v", t))
		}

		b.Tags = append(b.Tags, Tag{ID: t})
	}

	return nil
}

func (b *AdminBanner) ValidateAndUpdateTags(newTags []int) ([]int, error) {
	if len(newTags) == 0 {
		return nil, nil
	}

	if len(slices.Compact(newTags)) != len(newTags) {
		return nil, errors.New("list of tags_ids contain similar ids")
	}

	tagsInt := make([]int, 0, len(b.Tags))
	for _, t := range b.Tags {
		tagsInt = append(tagsInt, t.ID)
	}

	toDel := make([]int, 0)

	for _, t := range newTags {
		if t == 0 {
			continue
		}
		if slices.Contains(tagsInt, t) {
			b.Tags = slices.Delete(b.Tags, slices.Index(tagsInt, t), slices.Index(tagsInt, t)+1)
			tagsInt = slices.Delete(tagsInt, slices.Index(tagsInt, t), slices.Index(tagsInt, t)+1)
			toDel = append(toDel, t)
			continue
		}

		b.Tags = append(b.Tags, Tag{ID: t})
	}

	fmt.Println("tags: ", b.Tags)
	fmt.Println(toDel)

	return toDel, nil
}

func (b *AdminBanner) ValidateAndSetFeature(feature int) error {
	if feature == -1 {
		b.Feature.ID = 0
		return nil
	}

	if feature < 0 {
		return errors.New("feature_id must be greater than 0")
	}

	b.Feature.ID = feature

	return nil
}
