package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestAuthClient_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessClient := mock.NewMockAuthHandlerClient(ctrl)

	cl := &AuthClient{client: sessClient}

	sessClient.EXPECT().Create(context.Background(), models.ModelSessionToGrpc(&models.Session{})).Return(&emptypb.Empty{}, nil)
	err := cl.Create(&models.Session{})
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	sessClient.EXPECT().Create(context.Background(), models.ModelSessionToGrpc(&models.Session{})).Return(nil, grpc.ErrClientConnClosing)
	err = cl.Create(&models.Session{})
	assert.Equal(t, err, errors.GRPCError(grpc.ErrClientConnClosing))
}

func TestAuthClient_Get(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessClient := mock.NewMockAuthHandlerClient(ctrl)
	cl := &AuthClient{client: sessClient}

	sessClient.EXPECT().Get(context.Background(), &auth.SessionValue{Value: "qwerty"}).Return(&auth.Session{}, nil)
	_, err := cl.Get("qwerty")
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	sessClient.EXPECT().Get(context.Background(), &auth.SessionValue{Value: "qwerty"}).Return(nil, grpc.ErrClientConnClosing)
	_, err = cl.Get("qwerty")
	assert.Equal(t, err, errors.GRPCError(grpc.ErrClientConnClosing))
}

func TestAuthClient_Check(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessClient := mock.NewMockAuthHandlerClient(ctrl)
	cl := &AuthClient{client: sessClient}

	sessClient.EXPECT().Check(context.Background(), &auth.SessionValue{Value: "qwerty"}).Return(&auth.Session{}, nil)
	_, err := cl.Check("qwerty")
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	sessClient.EXPECT().Check(context.Background(), &auth.SessionValue{Value: "qwerty"}).Return(nil, grpc.ErrClientConnClosing)
	_, err = cl.Check("qwerty")
	assert.Equal(t, err, errors.GRPCError(grpc.ErrClientConnClosing))
}

func TestAuthClient_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessClient := mock.NewMockAuthHandlerClient(ctrl)
	cl := &AuthClient{client: sessClient}

	sessClient.EXPECT().Delete(context.Background(), &auth.SessionValue{Value: "qwerty"}).Return(&emptypb.Empty{}, nil)
	err := cl.Delete("qwerty")
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	sessClient.EXPECT().Delete(context.Background(), &auth.SessionValue{Value: "qwerty"}).Return(nil, grpc.ErrClientConnClosing)
	err = cl.Delete("qwerty")
	assert.Equal(t, err, errors.GRPCError(grpc.ErrClientConnClosing))
}