package main

type Article struct {
	Title     string
	URL       string
	Latitude  string
	Longitude string
	Details   map[string]string
}

func NewArticle(title, url, latitude, longitude string, details map[string]string) *Article {
	return &Article{
		Title:     title,
		URL:       url,
		Latitude:  latitude,
		Longitude: longitude,
		Details:   details,
	}
}
