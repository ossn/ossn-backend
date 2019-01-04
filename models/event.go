package models

import (
	"strconv"
	"time"
)

type Event struct {
	Model
	Title           string     `json:"title" sql:"not null"`
	StartDate       *time.Time `json:"startDate"`
	EndDate         *time.Time `json:"endDate"`
	Location        *Location  `json:"location" gorm:"foreignkey:LocationID;association_foreignkey:ID"`
	LocationID      *uint      `json:"location_id"`
	ImageURL        *string    `json:"imageUrl" sql:"type:text;"`
	Description     *string    `json:"description" sql:"type:text;"`
	SortDescription *string    `json:"sortDescription" sql:"type:text;"`
	Club            *Club      `json:"club" gorm:"foreignkey:ClubID;"`
	ClubID          *uint      `json:"club_id"`
	PublishedAt     *time.Time `json:"publishedAt" gorm:"index:event_published_at"`
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
