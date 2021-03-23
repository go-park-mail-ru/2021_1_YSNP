package session

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
)

type SessionUsecase interface {
	Create(sess *models.Session) *errors.Error
	Get(sessValue string) (*models.Session, *errors.Error)
	IsExist(sessionValue string) bool
	Delete(sessionValue string) *errors.Error
	Check(sessValue string) (*models.Session, *errors.Error)
}