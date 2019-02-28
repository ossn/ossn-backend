package models

import (
	"fmt"
	"os"

	"gopkg.in/go-playground/validator.v9"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/qor/admin"
)

var (
	DBSession     *gorm.DB
	AdminResource *admin.Admin
	dbURL         = os.Getenv("DATABASE_URL")
	dbPassword    = os.Getenv("DB_PASSWORD")
	host          = os.Getenv("DB_HOST")
	user          = os.Getenv("DB_USER")
	database      = os.Getenv("DB_NAME")
	validate      *validator.Validate
)

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

	validate = validator.New()

	err = validate.RegisterValidation("notblank", notBlankValidation)
	if err != nil {
		fmt.Println("Error on registering non-blank validation: " + err.Error())
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
	)

	AdminResource = admin.New(&admin.AdminConfig{DB: DBSession, SiteName: "OSSN Admin"})
	AdminResource.AddResource(&Event{})
	AdminResource.AddResource(&Announcement{})
	AdminResource.AddResource(&Location{})
	AdminResource.AddResource(&Job{})
	AdminResource.AddResource(&ClubUserRole{})
	AdminResource.AddResource(&Club{})
	AdminResource.AddResource(&User{})

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

	// seed()

}
