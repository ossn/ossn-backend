package models

type Club struct {
	Model
	Email           *string         `json:"email" gorm:"unique_index"`
	Location        *Location       `json:"location"`
	Title           *string         `json:"name"`
	ImageURL        *string         `json:"imageUrl" sql:"type:text;"`
	Description     *string         `json:"description" sql:"type:text;"`
	CodeOfConduct   *string         `json:"codeOfConduct" sql:"type:text;"`
	SortDescription *string         `json:"sortDescription" sql:"type:text;"`
	Users           []*ClubUserRole `json:"users" gorm:"many2many:club_user;"`
	Events          []*Event        `json:"events" gorm:"many2many:club_events;"`
	GithubURL       *string         `json:"githubUrl" sql:"type:text;" gorm:"unique_index"`
	ClubURL         *string         `json:"clubUrl" sql:"type:text;" gorm:"unique_index"`
}

type ClubWithRole struct {
	Club
	Role
}

func (c *Club) CreateClub() error {
	return nil
}
