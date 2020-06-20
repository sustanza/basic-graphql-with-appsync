package models

type Coffee struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Origin string `json:"origin"`
	Roast  string `json:"roast"`
}
