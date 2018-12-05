package models

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

type Model struct {
	*gorm.Model
}

func (m *Model) IDToString() (string, error) {
	return strconv.Itoa(m.ID), nil
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
