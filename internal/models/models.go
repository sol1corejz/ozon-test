package models

type URL struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	ShortUrl string `json:"short_url"`
}
