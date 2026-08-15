package models

import "time"

type URL struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	Original  string `json:"original"`
	CreatedAt time.Time `json:"created_at"`
}

type ShortenInput struct {
	URL string `json:"url" binding:"required, url"`
	//"url" validator checks its a valid URL formate automatically
}

type URLStats struct {
	Code string `json:"code"`
	Original string `json:"original"`
	Hits int64 `json:"hits"`
	Shorter string `json:"short_url"`
}