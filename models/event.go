package models

import (
	"context"
	"strconv"
	"time"
)

type Event struct {
	Model
	Title           string     `json:"title" sql:"not null"  validate:"notblank"`
	StartDate       *time.Time `json:"startDate"`
	EndDate         *time.Time `json:"endDate"`
	Location        *Location  `json:"location" gorm:"foreignkey:LocationID;association_foreignkey:ID"`
	LocationID      *uint      `json:"locationId"`
	ImageURL        *string    `json:"imageUrl" sql:"type:text;" mod:"trim"`
	Description     *string    `json:"description" sql:"type:text;"`
	SortDescription *string    `json:"sortDescription" sql:"type:text;"`
	Club            *Club      `json:"club" gorm:"foreignkey:ClubID;"`
	ClubID          *uint      `json:"clubId"`
	PublishedAt     *time.Time `json:"publishedAt" gorm:"index:event_published_at"`
}

func (e *Event) BeforeSave() error {
	err := transformer.Struct(context.Background(), e)
	if err != nil {
		return err
	}

	err = validateHttp(e.ImageURL, "Image url", true, true)
	if err != nil {
		return err
	}

	return transformValidationError(validate.Struct(e))
}

func (e *Event) StartDateToString() (*string, error) {
	if e.StartDate == nil {
		return nil, nil
	}
	str := strconv.FormatInt(e.StartDate.Unix(), 10)
	return &str, nil
}

func (e *Event) EndDateToString() (*string, error) {
	if e.EndDate == nil {
		return nil, nil
	}
	str := strconv.FormatInt(e.EndDate.Unix(), 10)
	return &str, nil
}

func (e *Event) PublishedAtToString() (*string, error) {
	if e.PublishedAt == nil {
		return nil, nil
	}
	str := strconv.FormatInt(e.PublishedAt.Unix(), 10)
	return &str, nil
}
