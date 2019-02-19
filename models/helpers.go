package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
)

type Model struct {
	gorm.Model
}

var rebuildURL = os.Getenv("REBUILD_URL")

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
	http.Post(rebuildURL, "application/json", bytes.NewReader(requestByte))
}

func (m *Model) AfterCreate(*gorm.Scope) error {
	go rebuildFrontEnd()
	return nil
}

func (m *Model) AfterUpdate(*gorm.Scope) error {
	go rebuildFrontEnd()
	return nil
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
		{Email: "test1@test.com", Name: "Test Test", UserName: "username", OIDCID: "username"},
		{Email: "kevinvnle@gmail.com", Name: "Kevin Viet Le", UserName: "le", PersonalURL: &le, OIDCID: "le"},
		{Name: "Shadi Nasser Moustafa", Email: "snasser2015@my.fit.edu", UserName: "mo", PersonalURL: &mo, OIDCID: "mo"},
		{UserName: "nelson.perezliveedpun", Name: "Nelson Perez", Email: "nelson.perez@live.edpuniversity.edu", PersonalURL: &ne, OIDCID: "nels"},
		{UserName: "dan", Name: "Daniel DAlonzo", Email: "founder@actionhorizon.institute", PersonalURL: &da, OIDCID: "dan"},
		{UserName: "ve", Name: "Veronica Armour", Email: "veronica.armour@shu.edu"},
		{Name: "Carla Rodriguez y", UserName: "CarlaRodriguezy Calderón", Email: "CarlaRodriguezy.Calderón@acm.com", OIDCID: "ve"},
		{Email: "test@test.com", Name: "Test Test", UserName: "username1", OIDCID: "user"},
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
		{ClubID: clubs[0].ID, UserID: users[0].ID, Role: "member"},
		{ClubID: clubs[0].ID, UserID: users[7].ID, Role: "club_owner"},
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
