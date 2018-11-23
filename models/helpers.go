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

func (m *Model) CreatedAtToString() (string, error) {
	return strconv.FormatInt(m.CreatedAt.Unix(), 10), nil
}

func (m *Model) UpdatedAtToString() (string, error) {
	return strconv.FormatInt(m.UpdatedAt.Unix(), 10), nil
}
func (m *Model) DeletedAtToString() (*string, error) {
	str := strconv.FormatInt(m.DeletedAt.Unix(), 10)
	return &str, nil
}

func (m *Model) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	scope.SetColumn("ID", uuid)
	return nil
}
