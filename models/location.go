package models

type Location struct {
	*Model
	Address *string `json:"address" sql:"type:text;"`
	Lat     *string `json:"lat"`
	Lng     *string `json:"lng"`
}
