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
