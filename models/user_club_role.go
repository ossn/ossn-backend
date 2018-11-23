package models

type ClubUserRole struct {
	Model
	User User `json:"user"`
	Club Club `json:"event"`
	Role Role `json:"role"`
}
