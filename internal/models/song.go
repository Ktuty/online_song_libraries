package models

type Songs struct {
	ID          int    `json:"id"`
	Song        string `json:"song"`
	Group       string `json:"group"`
	Text        string `json:"text"`
	ReleaseDate string `json:"releaseDate"`
	Link        string `json:"link"`
}
