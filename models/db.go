package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DBSession *gorm.DB

func init() {

	var err error
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	DBSession, err = gorm.Open("postgres", "postgres://"+user+":"+password+"@"+host+"/ossn_backend?sslmode=disable")
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
	)

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

func seed() {
	star := "test.com"
	c := &Club{ClubURL: &star}
	err := DBSession.Create(c).Error
	if err != nil {
		return
	}
	user := &User{Email: "test1@test.com", FirstName: "Test", LastName: "Test", Password: "test123", UserName: "username"}
	err = DBSession.Create(&Announcement{}).Error
	if err != nil {
		return
	}
	err = DBSession.Create(&Job{}).Error
	if err != nil {
		return
	}
	err = DBSession.Create(user).Error
	if err != nil {
		return
	}

	st := "an address"
	loc := &Location{Address: &st}
	err = DBSession.Create(loc).Error
	if err != nil {
		return
	}

	str := "test1.com"
	club := &Club{ClubURL: &str, Location: loc}
	err = DBSession.Create(club).Error
	if err != nil {
		return
	}
	err = DBSession.Create(&ClubUserRole{ClubID: club.ID, UserID: user.ID, Role: "user"}).Error
	if err != nil {
		return
	}

	user = &User{Email: "test@test.com", FirstName: "Test", LastName: "Test", Password: "test123", UserName: "username1"}
	err = DBSession.Create(user).Error
	if err != nil {
		return
	}

	err = DBSession.Create(&ClubUserRole{ClubID: club.ID, UserID: user.ID, Role: "admin"}).Error
	if err != nil {
		return
	}

	err = DBSession.Create(&Event{Title: "test event", ClubID: &club.ID}).Error
	if err != nil {
		return
	}
}
