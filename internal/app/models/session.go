package models

import "time"

type Session struct {
	ID        uint64
	Value     string
	UserID    uint64
	ExpiresAt time.Time
}