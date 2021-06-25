package main

type Article struct {
	Title     string
	URL       string
	Latitude  float64
	Longitude float64
	Read      string
	Text      string
	Details   map[string]string
}

func NewArticle(title, url, read, text string, latitude, longitude float64, details map[string]string) *Article {
	return &Article{
		Title:     title,
		URL:       url,
		Latitude:  latitude,
		Longitude: longitude,
		Read:      read,
		Text:      text,
		Details:   details,
	}
}
