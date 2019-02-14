package models

type User struct {
	Model
	Email             string          `json:"email" gorm:"UNIQUE;"`
	UserName          string          `json:"userName" gorm:"UNIQUE;" sql:"not null"`
	Name              string          `json:"lastName" sql:"not null"`
	ImageURL          *string         `json:"imageUrl" sql:"type:text;"`
	ReceiveNewsletter *bool           `json:"receiveNewsletter" gorm:"default:false;"`
	Description       *string         `json:"description" sql:"type:text;"`
	SortDescription   *string         `json:"sortDescription" sql:"type:text;"`
	Clubs             []*ClubUserRole `json:"clubs"`
	GithubURL         *string         `json:"githubUrl" sql:"type:text;"`
	PersonalURL       *string         `json:"personalUrl" sql:"type:text;"`
	OIDCID            string          `json:"oidcId" gorm:"index:idx_oidc_id;UNIQUE;column:oidc_id"`
	AccessToken       string          `json:"accessToken" gorm:"not null; index:idx_access_token"`
}
