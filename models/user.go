package models

type User struct {
	Model
	Email             string          `json:"email" gorm:"unique_index"`
	Password          string          `json:"password"`
	UserName          string          `json:"userName" gorm:"unique_index"`
	FirstName         string          `json:"firstName"`
	LastName          string          `json:"lastName"`
	ImageURL          *string         `json:"imageUrl" sql:"type:text;"`
	ReceiveNewsletter *bool           `json:"receiveNewsletter"`
	Description       *string         `json:"description" sql:"type:text;"`
	SortDescription   *string         `json:"sortDescription" sql:"type:text;"`
	Clubs             []*ClubUserRole `json:"clubs" gorm:"many2many:club_users;"`
	GithubURL         *string         `json:"githubUrl" sql:"type:text;"`
	PersonalURL       *string         `json:"personalUrl" sql:"type:text;"`
}

type UserWithRole struct {
	User
	Role
}
