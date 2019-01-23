package models

type Session struct {
	Model
	Token       string `json:"token" gorm:"index:idx_token"`
	Cookie      string `json:"cookie" gorm:"index:idx_cookie"`
	User        User   `json:"user" gorm:"foreignkey:UserID"`
	UserID      uint   `json:"userId"`
	AccessToken string `json:"accessToken" gorm:"not null; index:idx_access_token"`
}
