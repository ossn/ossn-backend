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

var (
	DBSession     *gorm.DB
	AdminResource *admin.Admin
	Auth          *auth.Auth
	dbURL         = os.Getenv("DATABASE_URL")
	dbPassword    = os.Getenv("DB_PASSWORD")
	host          = os.Getenv("DB_HOST")
	user          = os.Getenv("DB_USER")
	database      = os.Getenv("DB_NAME")
)

// hotfix for https://github.com/qor/auth/issues/16
type hotfixedAuthIdentity auth_identity.AuthIdentity

func (hotfixedAuthIdentity) TableName() string { return "basics" }

func init() {

	var err error
	if len(dbURL) < 1 {
		dbURL = "postgres://" + user + ":" + dbPassword + "@" + host + "/" + database + "?sslmode=disable"
	}

	DBSession, err = gorm.Open("postgres", dbURL)

	// DBSession.LogMode(true)

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
		&Session{},
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

	err = DBSession.Model(&ClubUserRole{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		fmt.Println("Error on foreign key creation: " + err.Error())
		//TODO: Fix this
		// panic(err)
	}
	err = DBSession.Model(&ClubUserRole{}).AddForeignKey("club_id", "clubs(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		fmt.Println("Error on foreign key creation: " + err.Error())
		//TODO: Fix this
		// panic(err)
	}
	err = DBSession.Model(&Event{}).AddForeignKey("location_id", "locations(id)", "RESTRICT", "CASCADE").Error
	if err != nil {
		fmt.Println("Error on foreign key creation: " + err.Error())
		//TODO: Fix this
		// panic(err)
	}
	err = DBSession.Model(&Event{}).AddForeignKey("club_id", "clubs(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		fmt.Println("Error on foreign key creation: " + err.Error())
		//TODO: Fix this
		// panic(err)
	}
	err = DBSession.Model(&Club{}).AddForeignKey("location_id", "locations(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		fmt.Println("Error on foreign key creation: " + err.Error())
		//TODO: Fix this
		// panic(err)
	}

	err = DBSession.Model(&Session{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		fmt.Println("Error on foreign key creation: " + err.Error())
		//TODO: Fix this
		// panic(err)
	}

	// seed()

}
