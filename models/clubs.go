package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

type Club struct {
	Model
	Email           *string   `json:"email" gorm:"UNIQUE;not null"  validate:"required,email"`
	Location        *Location `json:"location" gorm:"foreignkey:LocationID"`
	LocationID      *uint
	Title           *string         `json:"name" validate:"notblank" gorm:"UNIQUE;not null"`
	ImageURL        *string         `json:"imageUrl" sql:"type:text;"`
	Description     *string         `json:"description" sql:"type:text;"`
	CodeOfConduct   *string         `json:"codeOfConduct" sql:"type:text;"`
	SortDescription *string         `json:"sortDescription" sql:"type:text;"`
	Users           []*ClubUserRole `json:"users"`
	Events          []*Event        `json:"events"`
	GithubURL       *string         `json:"githubUrl" sql:"type:text;"`
	ClubURL         *string         `json:"clubUrl" sql:"type:text;"`
	BannerImageURL  *string         `json:"bannerImageUrl" sql:"type:text;"`
}

func (c *Club) AfterDelete(tx *gorm.DB) (err error) {
	return tx.Unscoped().Where("club_id = ?", c.ID).Delete(&ClubUserRole{}).Error
}

func (c *Club) BeforeSave() error {
	err := validate.Struct(c)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return errors.New(validationErrors.Error())
	}
	return nil
}
