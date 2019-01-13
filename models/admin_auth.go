package models

import (
	"fmt"

	"github.com/qor/admin"
	"github.com/qor/qor"
)

type AdminAuth struct{}

// var Auth = clean.New(&auth.Config{
// 	DB: DBSession,
// 	// User model needs to implement qor.CurrentUser interface (https://godoc.org/github.com/qor/qor#CurrentUser) to use it in QOR Admin
// 	UserModel: &Admin{},
// 	// AuthIdentityModel: &auth_identity.AuthIdentity{},
// })

func (AdminAuth) LoginURL(c *admin.Context) string {
	return "/auth/login"
}

func (AdminAuth) LogoutURL(c *admin.Context) string {
	return "/auth/logout"
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
