package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/qor/admin"
	"github.com/qor/auth"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/auth_themes/clean"
)

var DBSession *gorm.DB
var AdminResource *admin.Admin
var Auth *auth.Auth

// hotfix for https://github.com/qor/auth/issues/16
type hotfixedAuthIdentity auth_identity.AuthIdentity

func (hotfixedAuthIdentity) TableName() string { return "basics" }

func init() {

	var err error
	dbPassword := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	database := os.Getenv("DB_NAME")
	DBSession, err = gorm.Open("postgres", "postgres://"+user+":"+dbPassword+"@"+host+"/"+database+"?sslmode=disable")
	DBSession.LogMode(true)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	// Migrate the schema
	DBSession.Debug().AutoMigrate(
		&Event{},
		&Announcement{},
		&Location{},
		&Club{},
		&Job{},
		&ClubUserRole{},
		&User{},
		&Admin{},
	)

	Auth = clean.New(&auth.Config{
		DB:                DBSession,
		UserModel:         &Admin{},
		AuthIdentityModel: &hotfixedAuthIdentity{},
	})

	AdminResource = admin.New(&admin.AdminConfig{DB: DBSession, SiteName: "OSSN Admin", Auth: AdminAuth{}})
	AdminResource.AddResource(&Event{})
	AdminResource.AddResource(&Announcement{})
	AdminResource.AddResource(&Location{})
	AdminResource.AddResource(&Job{})
	AdminResource.AddResource(&ClubUserRole{}, &admin.Config{Invisible: true})
	AdminResource.AddResource(&Club{})
	AdminResource.AddResource(&User{})
	AdminResource.AddResource(&Admin{}, &admin.Config{Invisible: true})

	// .Meta(&admin.Meta{
	// 	Name: "Clubs",
	// 	// Valuer: func(record interface{}, c *qor.Context) interface{} {
	// 	// 	u, ok :=record.(ClubUserRole)
	// 	// 	fmt.Println(ok,u )
	// 	// 	// c.SetDB(c.GetDB().Preload("Club"))
	// 	// 	if ok && u.Club.Title != nil {
	// 	// 		return *u.Club.Title
	// 	// 	}
	// 	// 	return record
	// 	//  },
	// 	Config: &admin.SelectManyConfig{SelectOneConfig: admin.SelectOneConfig{RemoteDataResource: c}},
	// })

	err = DBSession.Model(&ClubUserRole{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		// panic(err)
	}
	err = DBSession.Model(&ClubUserRole{}).AddForeignKey("club_id", "clubs(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		// panic(err)
	}
	err = DBSession.Model(&Event{}).AddForeignKey("location_id", "locations(id)", "RESTRICT", "CASCADE").Error
	if err != nil {
		// panic(err)
	}
	err = DBSession.Model(&Event{}).AddForeignKey("club_id", "clubs(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		// panic(err)
	}
	err = DBSession.Model(&Club{}).AddForeignKey("location_id", "locations(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		// panic(err)
	}
	seed()

}
