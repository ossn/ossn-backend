package models

import (
	"context"
	"errors"

	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

type Club struct {
	Model
	Email           *string   `json:"email" gorm:"UNIQUE;not null" mod:"trim" validate:"required,email"`
	Location        *Location `json:"location" gorm:"foreignkey:LocationID"`
	LocationID      *uint
	Title           *string         `json:"name" validate:"required,notblank" gorm:"UNIQUE;not null"`
	ImageURL        *string         `json:"imageUrl" mod:"trim" sql:"type:text;"`
	Description     *string         `json:"description" sql:"type:text;"`
	CodeOfConduct   *string         `json:"codeOfConduct" sql:"type:text;"`
	SortDescription *string         `json:"sortDescription" sql:"type:text;"`
	Users           []*ClubUserRole `json:"users"`
	Events          []*Event        `json:"events"`
	GithubURL       *string         `json:"githubUrl" mod:"trim" sql:"type:text;"`
	ClubURL         *string         `json:"clubUrl" mod:"trim" sql:"type:text"`
	BannerImageURL  *string         `json:"bannerImageUrl" mod:"trim" sql:"type:text;"`
}

func (c *Club) AfterDelete(tx *gorm.DB) (err error) {
	return tx.Unscoped().Where("club_id = ?", c.ID).Delete(&ClubUserRole{}).Error
}

func (c *Club) BeforeSave() error {
	err := transformer.Struct(context.Background(), c)
	if err != nil {
		return err
	}

	err = validateHttp(c.GithubURL, "Github url", false, false)
	if err != nil {
		return err
	}

	err = validateHttp(c.ImageURL, "Image url", true, true)
	if err != nil {
		return err
	}

	err = validateHttp(c.BannerImageURL, "Banner image url", true, true)
	if err != nil {
		return err
	}

	err = validateHttp(c.ClubURL, "Club url", false, false)
	if err != nil {
		return err
	}

	err = validate.Struct(c)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return errors.New(validationErrors.Error())
	}
	return nil
}
