package usecase

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"time"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session"
)

type SessionUsecase struct {
	sessRepo session.SessionRepository
}

func NewSessionUsecase(repo session.SessionRepository) session.SessionUsecase {
	return &SessionUsecase{
		sessRepo: repo,
	}
}

func (su *SessionUsecase) Create(sess *models.Session) *errors.Error {
	err := su.sessRepo.Insert(sess)
	if err != nil {
		return errors.UnexpectedInternal(err)
	}

	return nil
}

func (su *SessionUsecase) Get(sessValue string) (*models.Session, *errors.Error) {
	sess, err := su.sessRepo.SelectByValue(sessValue)
	if err != nil {
		return nil, errors.Cause(errors.SessionNotExist)
	}

	return sess, nil
}

func (su *SessionUsecase) Delete(sessionValue string) *errors.Error {
	if _, err := su.Get(sessionValue); err != nil {
		return errors.Cause(errors.SessionNotExist)
	}

	err := su.sessRepo.DeleteByValue(sessionValue)
	if err != nil {
		return errors.UnexpectedInternal(err)
	}

	return nil
}

func (su *SessionUsecase) Check(sessValue string) (*models.Session, *errors.Error) {
	sess, err := su.Get(sessValue)
	if err != nil {
		return nil, errors.Cause(errors.SessionNotExist)
	}

	if sess.ExpiresAt.Before(time.Now()) {
		errE := su.Delete(sessValue)
		if errE != nil {
			return nil, errE
		}

		return nil, errors.Cause(errors.SessionExpired)
	}

	return sess, nil
}
