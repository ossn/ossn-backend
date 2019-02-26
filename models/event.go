package models

import (
	"errors"
	"strconv"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type Event struct {
	Model
	Title           string     `json:"title" sql:"not null"  validate:"notblank"`
	StartDate       *time.Time `json:"startDate"`
	EndDate         *time.Time `json:"endDate"`
	Location        *Location  `json:"location" gorm:"foreignkey:LocationID;association_foreignkey:ID"`
	LocationID      *uint      `json:"locationId"`
	ImageURL        *string    `json:"imageUrl" sql:"type:text;"`
	Description     *string    `json:"description" sql:"type:text;"`
	SortDescription *string    `json:"sortDescription" sql:"type:text;"`
	Club            *Club      `json:"club" gorm:"foreignkey:ClubID;"`
	ClubID          *uint      `json:"clubId"`
	PublishedAt     *time.Time `json:"publishedAt" gorm:"index:event_published_at"`
}

func (e *Event) BeforeSave() error {
	err := validate.Struct(e)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return errors.New(validationErrors.Error())
	}
	return nil
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
