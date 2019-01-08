package models

type Admin struct {
	Model
	Email    string `json:"email" gorm:"UNIQUE;"`
	Password string `json:"password" sql:"not null"`
	UserName string `json:"userName" gorm:"UNIQUE;" sql:"not null"`
}

func (a Admin) DisplayName() string {
	return a.UserName
}
