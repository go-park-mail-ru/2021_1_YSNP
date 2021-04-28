package usecase

import (
	errors2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
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

func (su *SessionUsecase) Create(sess *models.Session) *errors2.Error {
	err := su.sessRepo.Insert(sess)
	if err != nil {
		return errors2.UnexpectedInternal(err)
	}

	return nil
}

func (su *SessionUsecase) Get(sessValue string) (*models.Session, *errors2.Error) {
	sess, err := su.sessRepo.SelectByValue(sessValue)
	if err != nil {
		return nil, errors2.Cause(errors2.SessionNotExist)
	}

	return sess, nil
}

func (su *SessionUsecase) Delete(sessionValue string) *errors2.Error {
	if _, err := su.Get(sessionValue); err != nil {
		return errors2.Cause(errors2.SessionNotExist)
	}

	err := su.sessRepo.DeleteByValue(sessionValue)
	if err != nil {
		return errors2.UnexpectedInternal(err)
	}

	return nil
}

func (su *SessionUsecase) Check(sessValue string) (*models.Session, *errors2.Error) {
	sess, err := su.Get(sessValue)
	if err != nil {
		return nil, errors2.Cause(errors2.SessionNotExist)
	}

	if sess.ExpiresAt.Before(time.Now()) {
		errE := su.Delete(sessValue)
		if errE != nil {
			return nil, errE
		}

		return nil, errors2.Cause(errors2.SessionExpired)
	}

	return sess, nil
}
