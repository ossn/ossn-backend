package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

type Model struct {
	gorm.Model
}

func (m *Model) IDToString() (string, error) {
	return strconv.FormatUint(uint64(m.ID), 10), nil
}

func (m *Model) CreatedAtToString() string {
	return strconv.FormatInt(m.CreatedAt.Unix(), 10)
}

func (m *Model) UpdatedAtToString() string {
	return strconv.FormatInt(m.UpdatedAt.Unix(), 10)
}
func (m *Model) DeletedAtToString() *string {
	str := strconv.FormatInt(m.DeletedAt.Unix(), 10)
	return &str
}

func rebuildFrontEnd() {
	requestByte, err := json.Marshal(struct{}{})
	if err != nil {
		fmt.Println("Error while rebuilding frontend", err)
		return
	}
	http.Post("http://api.netlify.com/build_hooks/5c190a084e4723016c4b43af", "application/json", bytes.NewReader(requestByte))
}

func (m *Model) AfterCreate(scope *gorm.Scope) (err error) {
	go rebuildFrontEnd()
	return
}

func (m *Model) AfterUpdate(scope *gorm.Scope) (err error) {
	go rebuildFrontEnd()
	return
}

func TurnStringToRolename(name string) *RoleName {
	switch name {
	case "admin":
		str := RoleNameAdmin
		return &str

	case "member":
		str := RoleNameMember
		return &str

	case "club_owner":
		str := RoleNameClubOwner
		return &str
	case "guest":
		str := RoleNameGuest
		return &str
	default:
		return nil

	}
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
		{Email: "test1@test.com", FirstName: "Test", LastName: "Test", Password: "test123", UserName: "username"},
		{Email: "kevinvnle@gmail.com", FirstName: "Kevin Viet", LastName: "Le", Password: "test123", UserName: "le", PersonalURL: &le},
		{FirstName: "Shadi Nasser", LastName: "Moustafa", Email: "snasser2015@my.fit.edu", Password: "test123", UserName: "mo", PersonalURL: &mo},
		{Password: "test123", UserName: "nelson.perezliveedpun", FirstName: "Nelson", LastName: "Perez", Email: "nelson.perez@live.edpuniversity.edu", PersonalURL: &ne},
		{Password: "test123", UserName: "dan", FirstName: "Daniel", LastName: "DAlonzo", Email: "founder@actionhorizon.institute", PersonalURL: &da},
		{Password: "test123", UserName: "ve", FirstName: "Veronica", LastName: "Armour", Email: "veronica.armour@shu.edu"},
		{FirstName: "Carla Rodriguez y", UserName: "CarlaRodriguezy", LastName: "Calderón", Email: "CarlaRodriguezy.Calderón@acm.com", Password: "test123"},
		{Email: "test@test.com", FirstName: "Test", LastName: "Test", Password: "test123", UserName: "username1"},
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

	acm := "ACM Mines university of charleston"
	acmu := "acm.com"
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
		{ClubURL: &acmu, Title: &acm},
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
		{ClubID: clubs[0].ID, UserID: users[7].ID, Role: "admin"},
		{ClubID: clubs[4].ID, UserID: users[1].ID, Role: "member"},
		{ClubID: clubs[3].ID, UserID: users[2].ID, Role: "member"},
		{ClubID: clubs[5].ID, UserID: users[3].ID, Role: "member"},
		{ClubID: clubs[2].ID, UserID: users[4].ID, Role: "member"},
		{ClubID: clubs[1].ID, UserID: users[5].ID, Role: "member"},
		{ClubID: clubs[6].ID, UserID: users[6].ID, Role: "member"},
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
