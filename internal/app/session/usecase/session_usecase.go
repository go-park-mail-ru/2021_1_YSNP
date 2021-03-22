package usecase

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
)

type SessionUsecase struct {
	sessRepo session.SessionRepository
}

func NewSessionUsecase(repo session.SessionRepository) session.SessionUsecase{
	return &SessionUsecase{
			sessRepo: repo,
		}
}

func (su *SessionUsecase) Create(sess *models.Session) *errors.Error {
	panic("implement me")
}

func (su *SessionUsecase) Get(sessValue string) (*models.Session, *errors.Error) {
	panic("implement me")
}

func (su *SessionUsecase) IsExist(sessionValue string) bool {
	panic("implement me")
}

func (su *SessionUsecase) Delete(sessionValue string) *errors.Error {
	panic("implement me")
}

func (su *SessionUsecase) Check(sessValue string) (*models.Session, *errors.Error) {
	panic("implement me")
}