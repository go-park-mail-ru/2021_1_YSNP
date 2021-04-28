package auth

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

//go:generate mockgen -destination=./mocks/mock_session_repo.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session SessionRepository

type SessionRepository interface {
	Insert(session *models.Session) error
	SelectByValue(sessValue string) (*models.Session, error)
	DeleteByValue(sessionValue string) error
}
