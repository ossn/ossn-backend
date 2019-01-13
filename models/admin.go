package models

import (
	"fmt"
	"time"

	"github.com/qor/admin"
	"github.com/qor/qor"
)

type Admin struct {
	Model
	Email    string `json:"email" gorm:"UNIQUE;" form:"email"`
	Password string `json:"password" sql:"not null"`
	Name     string `json:"name" gorm:"UNIQUE;" sql:"not null" form:"name"`

	// Confirm
	ConfirmToken string
	Confirmed    bool

	// Recover
	RecoverToken       string
	RecoverTokenExpiry *time.Time
}

func (a Admin) DisplayName() string {
	return a.Email
}

func (Admin) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	currentUser := Auth.GetCurrentUser(c.Request)
	if currentUser != nil {
		qorCurrentUser, ok := currentUser.(qor.CurrentUser)
		if !ok {
			fmt.Printf("User %#v haven't implement qor.CurrentUser interface\n", currentUser)
		}
		return qorCurrentUser
	}
	return nil
}
