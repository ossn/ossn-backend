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
	database := os.Getenv("DB_NAME")
	DBSession, err = gorm.Open("postgres", "postgres://"+user+":"+password+"@"+host+"/"+database+"?sslmode=disable")
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
	err = DBSession.Create(&Announcement{}).Error
	if err != nil {
		return
	}
	err = DBSession.Create(&Job{}).Error
	if err != nil {
		return
	}
	le := "https://mozillians-dev.allizom.org/en-US/u/kevinvle/"
	mo := "https://mozillians.org/en-US/u/snasser2015/"
	ne := "nelson6855"
	da := "http://danieldalonzo.com/mozilla-learning-club-collaboration/"
	users := []User{
		{
			Email: "test1@test.com", FirstName: "Test", LastName: "Test", Password: "test123", UserName: "username",
		},
		{
			Email: "kevinvnle@gmail.com", FirstName: "Kevin Viet", LastName: "Le", Password: "test123", UserName: "le", PersonalURL: &le,
		},
		{
			FirstName: "Shadi Nasser", LastName: "Moustafa", Email: "snasser2015@my.fit.edu", Password: "test123", UserName: "mo", PersonalURL: &mo,
		},
		{
			Password: "test123", UserName: "nelson.perezliveedpun", FirstName: "Nelson", LastName: "Perez", Email: "nelson.perez@live.edpuniversity.edu", PersonalURL: &ne,
		},
		{
			Password: "test123", UserName: "dan", FirstName: "Daniel", LastName: "DAlonzo", Email: "founder@actionhorizon.institute", PersonalURL: &da,
		},
		{
			Password: "test123", UserName: "le", FirstName: "Veronica", LastName: "Armour", Email: "veronica.armour@shu.edu",
		},
		{
			Email: "test@test.com", FirstName: "Test", LastName: "Test", Password: "test123", UserName: "username1",
		},
	}

	err = DBSession.Create(&users).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	st := "an address"
	loc := &Location{Address: &st}
	err = DBSession.Create(loc).Error
	if err != nil {
		return
	}

	str := "test1.com"

	f := "https://www.ucsc.edu/"
	ft := "Fifikos"
	fa := "1156 High St, Santa Cruz, 95064, CA, USA"

	fouf := "Foufoutos"
	foufu := "fit.edu"
	foufa := "150 W University Blvd., 32901, Melbourne, Florida, USA"

	d := "Dedomena"
	du := "www.edpuniversity.edu"
	deda := "Betances # 49 PO Box 1674, 685, San Sebastian,Puerto Rico,USA"

	syn := "Syneffo"
	synu := "actionhorizon.institute"
	syna := "4 Old Forge Road,7930,Chester,NJ,USA"

	ok := "Okeanos"
	oka := "400 South Orange Avenue, 7079, South Orange, New Jersey, USA"
	oku := "www.shu.edu"

	clubs := []Club{{ClubURL: &str, Location: loc},
		{ClubURL: &f, Title: &ft, Location: &Location{Address: &fa}},
		{ClubURL: &foufu, Title: &fouf, Location: &Location{Address: &foufa}},
		{ClubURL: &du, Title: &d, Location: &Location{Address: &deda}},
		{ClubURL: &synu, Title: &syn, Location: &Location{Address: &syna}},
		{ClubURL: &oku, Title: &ok, Location: &Location{Address: &oka}},
	}
	err = DBSession.Create(&clubs).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	curs := []ClubUserRole{
		{ClubID: clubs[0].ID, UserID: users[0].ID, Role: "user"},
		{ClubID: clubs[0].ID, UserID: users[len(users)-1].ID, Role: "admin"},
		{ClubID: clubs[4].ID, UserID: users[1].ID, Role: "member"},
		{ClubID: clubs[3].ID, UserID: users[2].ID, Role: "member"},
		{ClubID: clubs[5].ID, UserID: users[3].ID, Role: "member"},
		{ClubID: clubs[2].ID, UserID: users[4].ID, Role: "member"},
		{ClubID: clubs[1].ID, UserID: users[5].ID, Role: "member"},
	}
	err = DBSession.Create(&curs).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	err = DBSession.Create(&Event{Title: "test event", ClubID: &clubs[0].ID}).Error
	if err != nil {
		return
	}
}
