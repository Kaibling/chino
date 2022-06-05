package models

type Movie struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
	Genre       string `json:"genre"`
	Favorite    bool   `json:"favorite"`
	Notified    bool   `json:"notified"`
}
