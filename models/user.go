package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	Model
	Email             string          `json:"email" gorm:"UNIQUE;" validate:"required,email"`
	UserName          string          `json:"userName" gorm:"UNIQUE;" sql:"not null" validate:"notblank"`
	Name              string          `json:"name" sql:"not null" validate:"notblank"`
	ImageURL          *string         `json:"imageUrl" sql:"type:text;"`
	ReceiveNewsletter *bool           `json:"receiveNewsletter" gorm:"default:false;"`
	Description       *string         `json:"description" sql:"type:text;"`
	SortDescription   *string         `json:"sortDescription" sql:"type:text;"`
	Clubs             []*ClubUserRole `json:"clubs"`
	GithubURL         *string         `json:"githubUrl" sql:"type:text;"`
	PersonalURL       *string         `json:"personalUrl" sql:"type:text;"`
	OIDCID            string          `json:"oidcId" gorm:"index:idx_oidc_id;UNIQUE;column:oidc_id"`
	AccessToken       string          `json:"accessToken" gorm:"not null; index:idx_access_token"`
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	return tx.Unscoped().Where("user_id = ?", u.ID).Delete(&ClubUserRole{}).Error
}

func (u *User) BeforeSave() error {
	err := validate.Struct(u)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return errors.New(validationErrors.Error())
	}
	return nil
}
