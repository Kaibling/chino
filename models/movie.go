package models

type Movie struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
	Genre       string `json:"genre"`
	Favorite    bool   `json:"favorite"`
	Notified    bool   `json:"notified"`
}

type String string

type MemoryCheck struct {
	Alloc      string `json:"alloc"`
	TotalAlloc string `json:"total_alloc"`
	Sys        string `json:"sys"`
	NumGC      string `json:"num_gc"`
}

type HealthCheck struct {
}
