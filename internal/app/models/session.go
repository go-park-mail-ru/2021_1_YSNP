package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Value     string
	UserID    uint64
	ExpiresAt time.Time
}

type LoginRequest struct {
	Telephone string `json:"telephone" valid:"phoneNumber"`
	Password  string `json:"password" valid:"password"`
}

func CreateSession(userID uint64) *Session {
	return &Session{
		Value:     uuid.New().String(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(10 * time.Hour),
	}
}
