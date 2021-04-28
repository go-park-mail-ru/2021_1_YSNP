package session

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	errors2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

//go:generate mockgen -destination=./mocks/mock_session_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session SessionUsecase

type SessionUsecase interface {
	Create(sess *models.Session) *errors2.Error
	Get(sessValue string) (*models.Session, *errors2.Error)
	Delete(sessionValue string) *errors2.Error
	Check(sessValue string) (*models.Session, *errors2.Error)
}
