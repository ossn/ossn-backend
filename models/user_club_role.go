package models

import (
	"errors"
)

type ClubUserRole struct {
	Model
	UserID uint   `sql:"not null"`
	User   User   `json:"user" gorm:"foreignkey:UserID"`
	ClubID uint   `sql:"not null"`
	Club   Club   `json:"club" gorm:"foreignkey:ClubID;"`
	Role   string `json:"role" gorm:"default:'guest'" sql:"not null"`
}

func (c *ClubUserRole) BeforeSave() error {
	role := TurnStringToRolename(c.Role)
	if role == nil {
		return errors.New("Role can only be member, admin, club_owner or guest")
	}
	return nil
}
