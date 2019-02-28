package models

import (
	"context"
	"strconv"
	"time"
)

type Job struct {
	Model
	Description     *string    `json:"description" sql:"type:text;"`
	SortDescription *string    `json:"sortDescription" sql:"type:text;"`
	URL             *string    `json:"url" sql:"type:text;" mod:"trim"`
	ImageURL        *string    `json:"imageUrl" sql:"type:text;" mod:"trim"`
	PublishedAt     *time.Time `json:"publishedAt" gorm:"index:job_published_at"`
}

func (j *Job) BeforeSave() error {
	err := transformer.Struct(context.Background(), j)
	if err != nil {
		return err
	}

	err = validateHttp(j.ImageURL, "Image url", true, true)
	if err != nil {
		return err
	}

	err = validateHttp(j.URL, "Url", false, false)
	if err != nil {
		return err
	}

	return nil
}

func (j *Job) PublishedAtToString() (*string, error) {
	if j.PublishedAt == nil {
		return nil, nil
	}
	str := strconv.FormatInt(j.PublishedAt.Unix(), 10)
	return &str, nil
}
