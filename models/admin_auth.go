package models

import (
	"fmt"

	"github.com/qor/admin"
	"github.com/qor/auth"
	"github.com/qor/auth_themes/clean"
	"github.com/qor/qor"
)

type AdminAuth struct{}

var Auth = clean.New(&auth.Config{
	DB: DBSession,
	// User model needs to implement qor.CurrentUser interface (https://godoc.org/github.com/qor/qor#CurrentUser) to use it in QOR Admin
	UserModel: Admin{},
})

func (AdminAuth) LoginURL(c *admin.Context) string {
	return "/admin/auth/login"
}

func (AdminAuth) LogoutURL(c *admin.Context) string {
	return "/admin/auth/logout"
}

func (AdminAuth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
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
