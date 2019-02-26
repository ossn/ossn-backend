package models

import "github.com/jinzhu/gorm"

type Club struct {
	Model
	Email           *string   `json:"email" gorm:"UNIQUE;not null"`
	Location        *Location `json:"location" gorm:"foreignkey:LocationID"`
	LocationID      *uint
	Title           *string         `json:"name"`
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
