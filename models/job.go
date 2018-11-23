package models

import (
	"strconv"
	"time"
)

type Job struct {
	Model
	Description     *string    `json:"description" sql:"type:text;"`
	SortDescription *string    `json:"sortDescription" sql:"type:text;"`
	URL             *string    `json:"url" sql:"type:text;"`
	ImageURL        *string    `json:"imageUrl" sql:"type:text;"`
	PublishedAt     *time.Time `json:"publishedAt" gorm:"index:job_published_at"`
}

func (j *Job) PublishedAtToString() (*string, error) {
	if j.PublishedAt == nil {
		return nil, nil
	}
	str := strconv.FormatInt(j.PublishedAt.Unix(), 10)
	return &str, nil
}
