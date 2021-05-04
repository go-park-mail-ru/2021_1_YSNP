package models

import (
	proto "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/auth"
	"github.com/golang/protobuf/ptypes"
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

func GrpcSessionToModel(grpcSess *proto.Session) *Session {
	ExpiresAt, _ := ptypes.Timestamp(grpcSess.ExpiresAt)

	return &Session{
		Value:     grpcSess.Value,
		UserID:    grpcSess.UserID,
		ExpiresAt: ExpiresAt,
	}
}

func ModelSessionToGrpc(modelSess *Session) *proto.Session {
	ExpiresAt, _ := ptypes.TimestampProto(modelSess.ExpiresAt)

	return &proto.Session{
		Value:     modelSess.Value,
		UserID:    modelSess.UserID,
		ExpiresAt: ExpiresAt,
	}
}
