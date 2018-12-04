package models

import (
	"github.com/satori/go.uuid"
)

type ClubUserRole struct {
	Model
	UserID uuid.UUID `sql:"not null" gorm:"type:uuid"`
	User   User      `json:"user" gorm:"foreignkey:ID;association_foreignkey:UserID"`
	ClubID uuid.UUID `sql:"not null" gorm:"type:uuid"`
	Club   Club      `json:"club" gorm:"foreignkey:ID;association_foreignkey:ClubID;"`
	Role   string    `json:"role" gorm:"default:'guest'" sql:"not null"`
}
