package models

type User struct {
	Model
	Email             string          `json:"email" gorm:"UNIQUE;"`
	Password          string          `json:"password" sql:"not null"`
	UserName          string          `json:"userName" gorm:"UNIQUE;" sql:"not null"`
	FirstName         string          `json:"firstName" sql:"not null"`
	LastName          string          `json:"lastName" sql:"not null"`
	ImageURL          *string         `json:"imageUrl" sql:"type:text;"`
	ReceiveNewsletter *bool           `json:"receiveNewsletter" gorm:"default:false;"`
	Description       *string         `json:"description" sql:"type:text;"`
	SortDescription   *string         `json:"sortDescription" sql:"type:text;"`
	Clubs             []*ClubUserRole `json:"clubs" gorm:"many2many:club_user_roles"`
	GithubURL         *string         `json:"githubUrl" sql:"type:text;"`
	PersonalURL       *string         `json:"personalUrl" sql:"type:text;"`
}
