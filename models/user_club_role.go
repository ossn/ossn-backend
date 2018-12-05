package models

type ClubUserRole struct {
	*Model
	UserID uint   `sql:"not null"`
	User   User   `json:"user" gorm:"foreignkey:ID;association_foreignkey:UserID"`
	ClubID uint   `sql:"not null"`
	Club   Club   `json:"club" gorm:"foreignkey:ID;association_foreignkey:ClubID;"`
	Role   string `json:"role" gorm:"default:'guest'" sql:"not null"`
}
