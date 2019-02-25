// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"fmt"
	"io"
	"strconv"
)

type Announcements struct {
	Announcements []Announcement `json:"announcements"`
	PageInfo      PageInfo       `json:"pageInfo"`
}

func (Announcements) IsWithPagination() {}

type ClubInput struct {
	Email           string         `json:"email"`
	Location        *LocationInput `json:"location"`
	Name            string         `json:"name"`
	ImageURL        *string        `json:"imageUrl"`
	BannerImageURL  *string        `json:"bannerImageUrl"`
	Description     string         `json:"description"`
	CodeOfConduct   *string        `json:"codeOfConduct"`
	SortDescription *string        `json:"sortDescription"`
	GithubURL       *string        `json:"githubUrl"`
	ClubURL         *string        `json:"clubUrl"`
}

type ClubWithRole struct {
	ID              string    `json:"id"`
	Email           *string   `json:"email"`
	Location        *Location `json:"location"`
	Name            *string   `json:"name"`
	ImageURL        *string   `json:"imageUrl"`
	BannerImageURL  *string   `json:"bannerImageUrl"`
	Description     *string   `json:"description"`
	CodeOfConduct   *string   `json:"codeOfConduct"`
	SortDescription *string   `json:"sortDescription"`
	Users           []*User   `json:"users"`
	Events          []*Event  `json:"events"`
	Role            *RoleName `json:"role"`
	GithubURL       *string   `json:"githubUrl"`
	ClubURL         *string   `json:"clubUrl"`
	CreatedAt       string    `json:"createdAt"`
	UpdatedAt       string    `json:"updatedAt"`
}

type Clubs struct {
	Clubs    []Club   `json:"clubs"`
	PageInfo PageInfo `json:"pageInfo"`
}

func (Clubs) IsWithPagination() {}

type EventInput struct {
	Title           string  `json:"title"`
	StartDate       *string `json:"startDate"`
	EndDate         *string `json:"endDate"`
	LocationID      *string `json:"locationID"`
	ImageURL        *string `json:"imageUrl"`
	Description     *string `json:"description"`
	SortDescription *string `json:"sortDescription"`
	ClubID          string  `json:"clubId"`
}

type Events struct {
	Events   []Event  `json:"events"`
	PageInfo PageInfo `json:"pageInfo"`
}

func (Events) IsWithPagination() {}

type Jobs struct {
	Jobs     []Job    `json:"jobs"`
	PageInfo PageInfo `json:"pageInfo"`
}

func (Jobs) IsWithPagination() {}

type LocationInput struct {
	ID      *string `json:"id"`
	Address *string `json:"address"`
	Lat     *string `json:"lat"`
	Lng     *string `json:"lng"`
}

type PageInfo struct {
	StartCursor     string `json:"startCursor"`
	EndCursor       string `json:"endCursor"`
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	TotalCount      int    `json:"totalCount"`
}

type UserInput struct {
	Name              string   `json:"name"`
	ReceiveNewsletter bool     `json:"receiveNewsletter"`
	Description       *string  `json:"description"`
	SortDescription   *string  `json:"sortDescription"`
	Clubs             []string `json:"clubs"`
	GithubURL         *string  `json:"githubUrl"`
	PersonalURL       *string  `json:"personalUrl"`
}

type UserWithRole struct {
	ID                string          `json:"id"`
	Email             string          `json:"email"`
	UserName          string          `json:"userName"`
	Name              string          `json:"name"`
	ImageURL          *string         `json:"imageUrl"`
	ReceiveNewsletter *bool           `json:"receiveNewsletter"`
	Description       *string         `json:"description"`
	SortDescription   *string         `json:"sortDescription"`
	Clubs             []*ClubWithRole `json:"clubs"`
	GithubURL         *string         `json:"githubUrl"`
	PersonalURL       *string         `json:"personalUrl"`
	CreatedAt         string          `json:"createdAt"`
	UpdatedAt         string          `json:"updatedAt"`
	Role              *RoleName       `json:"role"`
}

type Users struct {
	Users    []User   `json:"users"`
	PageInfo PageInfo `json:"pageInfo"`
}

func (Users) IsWithPagination() {}

type WithPagination interface {
	IsWithPagination()
}

type RoleName string

const (
	RoleNameAdmin     RoleName = "admin"
	RoleNameMember    RoleName = "member"
	RoleNameClubOwner RoleName = "club_owner"
	RoleNameGuest     RoleName = "guest"
	RoleNameUndefined RoleName = "undefined"
)

func (e RoleName) IsValid() bool {
	switch e {
	case RoleNameAdmin, RoleNameMember, RoleNameClubOwner, RoleNameGuest, RoleNameUndefined:
		return true
	}
	return false
}

func (e RoleName) String() string {
	return string(e)
}

func (e *RoleName) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RoleName(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RoleName", str)
	}
	return nil
}

func (e RoleName) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
