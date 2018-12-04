package models

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type Model struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (m *Model) IDToString() (string, error) {
	return m.ID.String(), nil
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

func (m *Model) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	scope.SetColumn("ID", uuid)
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
