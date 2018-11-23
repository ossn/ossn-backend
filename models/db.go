package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DBSession *gorm.DB

func init() {

	var err error
	DBSession, err = gorm.Open("sqlite3", "test.db")
	// password := os.Getenv("DB_PASSWORD")
	// host := os.Getenv("DB_HOST")
	// user := os.Getenv("DB_USER")
	// DBSession, err = gorm.Open("postgres", "postgres://"+user+":"+password+"@"+host+"/ossn_backend?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	// Migrate the schema
	DBSession.Raw("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	DBSession.AutoMigrate(
		&Event{},
		&Announcement{},
		&Location{},
		&Club{},
		&Job{},
		&Role{},
		&User{},
	)

	user := &User{Email: "test@test.com", FirstName: "Test", LastName: "Test", Password: "test123"}
	DBSession.Create(user)
	str := "test.com"
	club := &Club{ClubURL: &str}
	DBSession.Create(club)
	DBSession.Model(user).Association("Clubs").Append(club)
	// DBSession.Create(&User{Clubs: })
	DBSession.Create(&Job{})
	DBSession.Create(&Event{Title: "test event"})
	DBSession.Create(&Announcement{})
}
