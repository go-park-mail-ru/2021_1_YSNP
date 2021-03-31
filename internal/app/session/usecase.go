package session

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

type SessionUsecase interface {
	Create(sess *models.Session) *errors.Error
	Get(sessValue string) (*models.Session, *errors.Error)
	Delete(sessionValue string) *errors.Error
	Check(sessValue string) (*models.Session, *errors.Error)
}
