package models

import "time"

type URL struct {
	ShortCode    string    `bson:"shortCode"`
	OriginalURL  string    `bson:"originalUrl"`
	CreatedAt    time.Time `bson:"createdAt"`
	ExpiresAt    time.Time `bson:"expiresAt"`
	Clicks       int       `bson:"clicks"`
	LastAccessed time.Time `bson:"lastAccessed,omitempty"`
}

type RateLimit struct {
	IP        string    `bson:"_id"`
	Remaining int       `bson:"remaining"`
	Reset     time.Time `bson:"reset"`
}

type Stats struct {
	TotalClicks int `bson:"value"`
}
