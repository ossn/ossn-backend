package models

import (
	"strconv"
	"time"
)

type Announcement struct {
	Model
	Description     *string    `json:"description" sql:"type:text;"`
	SortDescription *string    `json:"sortDescription" sql:"type:text;"`
	URL             *string    `json:"url" sql:"type:text;"`
	ImageURL        *string    `json:"imageUrl" sql:"type:text;"`
	PublishedAt     *time.Time `json:"publishedAt" gorm:"index:announcement_published_at"`
}

func (a *Announcement) BeforeSave() error {
	err := validateHttp(a.ImageURL, "Image url", true, true)
	if err != nil {
		return err
	}

	err = validateHttp(a.URL, "url", false, false)
	if err != nil {
		return err
	}
	return nil
}

func (a *Announcement) PublishedAtToString() (*string, error) {
	if a.PublishedAt == nil {
		return nil, nil
	}
	str := strconv.FormatInt(a.PublishedAt.Unix(), 10)
	return &str, nil
}
