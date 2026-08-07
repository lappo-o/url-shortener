package model

import "time"

type Response struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type Url struct {
	OriginalURL string
	ShortCode   string
	CreatedAt   time.Time
}
