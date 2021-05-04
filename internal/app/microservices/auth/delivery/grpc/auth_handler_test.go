package grpc

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/auth/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	proto "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/auth"
)

func TestAuthHandlerServer_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessUcase := mock.NewMockSessionUsecase(ctrl)

	sessHandler := NewAuthHandlerServer(sessUcase)

	sessUcase.EXPECT().Create(models.GrpcSessionToModel(&proto.Session{})).Return(nil)

	_, err := sessHandler.Create(context.Background(), &proto.Session{})
	assert.Equal(t, err, nil)

	//error
	sessUcase.EXPECT().Create(models.GrpcSessionToModel(&proto.Session{})).Return(errors.UnexpectedInternal(sql.ErrConnDone))

	_, err = sessHandler.Create(context.Background(), &proto.Session{})
	assert.Error(t, err)
}

func TestAuthHandlerServer_Get(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessUcase := mock.NewMockSessionUsecase(ctrl)

	sessHandler := NewAuthHandlerServer(sessUcase)

	sessUcase.EXPECT().Get("").Return(&models.Session{}, nil)

	_, err := sessHandler.Get(context.Background(), &proto.SessionValue{})
	assert.Equal(t, err, nil)

	//error
	sessUcase.EXPECT().Get("").Return(nil, errors.UnexpectedInternal(sql.ErrConnDone))

	_, err = sessHandler.Get(context.Background(), &proto.SessionValue{})
	assert.Error(t, err)
}

func TestAuthHandlerServer_Check(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessUcase := mock.NewMockSessionUsecase(ctrl)

	sessHandler := NewAuthHandlerServer(sessUcase)

	sessUcase.EXPECT().Check("").Return(&models.Session{}, nil)

	_, err := sessHandler.Check(context.Background(), &proto.SessionValue{})
	assert.Equal(t, err, nil)

	//error
	sessUcase.EXPECT().Check("").Return(nil, errors.UnexpectedInternal(sql.ErrConnDone))

	_, err = sessHandler.Check(context.Background(), &proto.SessionValue{})
	assert.Error(t, err)
}

func TestAuthHandlerServer_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessUcase := mock.NewMockSessionUsecase(ctrl)

	sessHandler := NewAuthHandlerServer(sessUcase)

	sessUcase.EXPECT().Delete("").Return(nil)

	_, err := sessHandler.Delete(context.Background(), &proto.SessionValue{})
	assert.Equal(t, err, nil)

	//error
	sessUcase.EXPECT().Delete("").Return(errors.UnexpectedInternal(sql.ErrConnDone))

	_, err = sessHandler.Delete(context.Background(), &proto.SessionValue{})
	assert.Error(t, err)
}
