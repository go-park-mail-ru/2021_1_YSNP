package usecase

import (
	errs "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSessionUsecase_Create_OK(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRep := mock.NewMockSessionRepository(ctrl)
	sessionUcase := NewSessionUsecase(sessionRep)

	session := models.CreateSession(0)

	sessionRep.EXPECT().Insert(gomock.Eq(session)).Return(nil)

	err := sessionUcase.Create(session)
	assert.Equal(t, err, (*errs.Error)(nil))
}

func TestSessionUsecase_Delete_OK(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRep := mock.NewMockSessionRepository(ctrl)
	sessionUcase := NewSessionUsecase(sessionRep)

	session := models.CreateSession(0)

	sessionRep.EXPECT().SelectByValue(session.Value).Return(session, nil)
	sessionRep.EXPECT().DeleteByValue(session.Value).Return(nil)

	err := sessionUcase.Delete(session.Value)
	assert.Equal(t, err, (*errs.Error)(nil))
}

func TestSessionUsecase_Delete_SessionNotExist(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRep := mock.NewMockSessionRepository(ctrl)
	sessionUcase := NewSessionUsecase(sessionRep)

	session := models.CreateSession(0)

	sessionRep.EXPECT().SelectByValue(session.Value).Return(nil, errors.New("session not exist"))

	err := sessionUcase.Delete(session.Value)
	assert.Equal(t, err, errs.Cause(errs.SessionNotExist))
}

func TestSessionUsecase_Get_OK(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRep := mock.NewMockSessionRepository(ctrl)
	sessionUcase := NewSessionUsecase(sessionRep)

	session := models.CreateSession(0)

	sessionRep.EXPECT().SelectByValue(session.Value).Return(session, nil)

	sess, err := sessionUcase.Get(session.Value)

	assert.Equal(t, err, (*errs.Error)(nil))
	assert.Equal(t, sess, session)
}

func TestSessionUsecase_Get_SessionNotExist(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRep := mock.NewMockSessionRepository(ctrl)
	sessionUcase := NewSessionUsecase(sessionRep)

	session := models.CreateSession(0)

	sessionRep.EXPECT().SelectByValue(session.Value).Return(nil, errors.New("cannot cast to string"))

	sess, err := sessionUcase.Get(session.Value)

	assert.Equal(t, err, errs.Cause(errs.SessionNotExist))
	assert.Equal(t, sess, (*models.Session)(nil))
}

func TestSessionUsecase_Check_OK(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRep := mock.NewMockSessionRepository(ctrl)
	sessionUcase := NewSessionUsecase(sessionRep)

	session := models.CreateSession(0)

	sessionRep.EXPECT().SelectByValue(session.Value).Return(session, nil)

	sess, err := sessionUcase.Check(session.Value)

	assert.Equal(t, err, (*errs.Error)(nil))
	assert.Equal(t, sess, session)
}

func TestSessionUsecase_Check_SessionNotExist(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRep := mock.NewMockSessionRepository(ctrl)
	sessionUcase := NewSessionUsecase(sessionRep)

	session := models.CreateSession(0)

	sessionRep.EXPECT().SelectByValue(session.Value).Return(nil, errors.New("session not exist"))

	_, err := sessionUcase.Check(session.Value)

	assert.Equal(t, err, errs.Cause(errs.SessionNotExist))
}

func TestSessionUsecase_Check_Expired(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRep := mock.NewMockSessionRepository(ctrl)
	sessionUcase := NewSessionUsecase(sessionRep)

	session := models.CreateSession(0)
	session.ExpiresAt = time.Now().AddDate(0, 0, -1)

	sessionRep.EXPECT().SelectByValue(session.Value).Return(session, nil).Times(2)
	sessionRep.EXPECT().DeleteByValue(session.Value).Return(nil)

	_, err := sessionUcase.Check(session.Value)

	assert.Equal(t, err, errs.Cause(errs.SessionExpired))
}
