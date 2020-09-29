package es

import "time"

type News struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
