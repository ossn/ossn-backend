package models

import (
	"context"

	"github.com/jinzhu/gorm"
)

type User struct {
	Model
	Email               string          `json:"email" gorm:"UNIQUE;" validate:"required,email"`
	UserName            string          `json:"userName" gorm:"UNIQUE;" sql:"not null" validate:"notblank"`
	Name                string          `json:"name" sql:"not null" validate:"notblank"`
	ImageURL            *string         `json:"imageUrl" sql:"type:text;"`
	ReceiveNewsletter   *bool           `json:"receiveNewsletter" gorm:"default:false;"`
	Description         *string         `json:"description" sql:"type:text;"`
	SortDescription     *string         `json:"sortDescription" sql:"type:text;"`
	Clubs               []*ClubUserRole `json:"clubs"`
	GithubURL           *string         `json:"githubUrl" sql:"type:text;"`
	PersonalURL         *string         `json:"personalUrl" sql:"type:text;"`
	OIDCID              string          `json:"oidcId" gorm:"index:idx_oidc_id;UNIQUE;column:oidc_id"`
	AccessToken         string          `json:"accessToken" gorm:"not null; index:idx_access_token"`
	IsOverTheLegalLimit bool            `json:"isOverTheLegalLimit" gorm:"default:false"`
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	return tx.Unscoped().Where("user_id = ?", u.ID).Delete(&ClubUserRole{}).Error
}

func (u *User) BeforeSave() error {
	err := transformer.Struct(context.Background(), u)
	if err != nil {
		return err
	}

	err = validateHttp(u.ImageURL, "Image url", false, false)
	if err != nil {
		return err
	}

	err = validateHttp(u.GithubURL, "Github url", false, false)
	if err != nil {
		return err
	}

	err = validateHttp(u.PersonalURL, "Personal url", false, false)
	if err != nil {
		return err
	}

	return transformValidationError(validate.Struct(u))
}
