package models

import "github.com/satori/go.uuid"

type Club struct {
	Model
	Email           *string         `json:"email" gorm:"UNIQUE;"`
	Location        *Location       `json:"location" gorm:"foreignkey:LocationID"`
	LocationID      *uuid.UUID      `gorm:"type:uuid"`
	Title           *string         `json:"name"`
	ImageURL        *string         `json:"imageUrl" sql:"type:text;"`
	Description     *string         `json:"description" sql:"type:text;"`
	CodeOfConduct   *string         `json:"codeOfConduct" sql:"type:text;"`
	SortDescription *string         `json:"sortDescription" sql:"type:text;"`
	Users           []*ClubUserRole `json:"users" gorm:"many2many:club_user_roles;association_foreignkey:ID;foreignkey:ID;association_jointable_foreignkey:user_id;jointable_foreignkey:club_id"`
	Events          []*Event        `json:"events" gorm:"many2many:club_events;"`
	GithubURL       *string         `json:"githubUrl" sql:"type:text;" gorm:"UNIQUE;"`
	ClubURL         *string         `json:"clubUrl" sql:"type:text;" gorm:"UNIQUE;"`
}
