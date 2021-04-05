package session

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

type SessionRepository interface {
	Insert(session *models.Session) error
	SelectByValue(sessValue string) (*models.Session, error)
	DeleteByValue(sessionValue string) error
}
