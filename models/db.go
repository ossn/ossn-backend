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
	tx := DBSession.Begin()

	err := tx.Create(&Announcement{}).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Create(&Job{}).Error
	if err != nil {
		tx.Rollback()
		return
	}
	le := "https://mozillians-dev.allizom.org/en-US/u/kevinvle/"
	mo := "https://mozillians.org/en-US/u/snasser2015/"
	ne := "nelson6855"
	da := "http://danieldalonzo.com/mozilla-learning-club-collaboration/"
	users := []User{
		User{
			Email: "test1@test.com", FirstName: "Test", LastName: "Test", Password: "test123", UserName: "username",
		},
		User{
			Email: "kevinvnle@gmail.com", FirstName: "Kevin Viet", LastName: "Le", Password: "test123", UserName: "le", PersonalURL: &le,
		},
		User{
			FirstName: "Shadi Nasser", LastName: "Moustafa", Email: "snasser2015@my.fit.edu", Password: "test123", UserName: "mo", PersonalURL: &mo,
		},
		User{
			Password: "test123", UserName: "nelson.perezliveedpun", FirstName: "Nelson", LastName: "Perez", Email: "nelson.perez@live.edpuniversity.edu", PersonalURL: &ne,
		},
		User{
			Password: "test123", UserName: "dan", FirstName: "Daniel", LastName: "DAlonzo", Email: "founder@actionhorizon.institute", PersonalURL: &da,
		},
		User{
			Password: "test123", UserName: "ve", FirstName: "Veronica", LastName: "Armour", Email: "veronica.armour@shu.edu",
		},
		User{
			Email: "test@test.com", FirstName: "Test", LastName: "Test", Password: "test123", UserName: "username1",
		},
	}
	for i, u := range users {

		err = tx.Create(&u).Error
		if err != nil {
			tx.Rollback()
			fmt.Println(u)
			return
		}
		users[i] = u
	}

	stx := "test1.com"

	star := "test.com"
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
	st := "an address"
	loc := []Location{{Address: &st}, {Address: &fa}, {Address: &foufa}, {Address: &deda}, {Address: &syna}, {Address: &oka}}
	for i, l := range loc {
		err = tx.Create(&l).Error
		if err != nil {
			tx.Rollback()
			return
		}
		loc[i] = l

	}

	clubs := []Club{
		{ClubURL: &stx, LocationID: &loc[0].ID},
		{ClubURL: &f, Title: &ft, LocationID: &loc[1].ID},
		{ClubURL: &foufu, Title: &fouf, LocationID: &loc[2].ID},
		{ClubURL: &du, Title: &d, LocationID: &loc[3].ID},
		{ClubURL: &synu, Title: &syn, LocationID: &loc[4].ID},
		{ClubURL: &oku, Title: &ok, LocationID: &loc[5].ID},
		{ClubURL: &star},
	}
		for i, c := range clubs {
		err = tx.Create(&c).Error
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			return
		}
		clubs[i] = c

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

		for i, c := range curs {
		err = tx.Create(&c).Error
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			return
		}
		curs[i] = c

	}


	err = tx.Create(&Event{Title: "test event", ClubID: &clubs[0].ID}).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
}
