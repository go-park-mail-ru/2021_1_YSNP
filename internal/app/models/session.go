package models

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID        uint64
	Value     string
	UserID    uint64
	ExpiresAt time.Time
}

type LoginRequest struct {
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

func CreateSession(userID uint64) *Session {
	return &Session{
		Value:     uuid.New().String(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(10 * time.Hour),
	}
}