package models

type Session struct {
	Token  string `json:"token" gorm:"index:idx_token"`
	UserID uint   `json:"userId"`
}
