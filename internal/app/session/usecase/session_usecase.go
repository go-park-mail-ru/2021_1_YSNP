package usecase

import (
	"context"
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/auth"
)

type AuthClient struct {
	client auth.AuthHandlerClient
}

func NewAuthClient(conn *grpc.ClientConn) *AuthClient {
	c := auth.NewAuthHandlerClient(conn)
	return &AuthClient{
		client: c,
	}
}

func (ac *AuthClient) Create(sess *models.Session) *errors.Error {
	_, err := ac.client.Create(context.Background(), models.ModelSessionToGrpc(sess))
	if err != nil {
		return errors.GRPCError(err)
	}
	return nil
}

func (ac *AuthClient) Get(sessVal string) (*models.Session, *errors.Error) {
	sess, err := ac.client.Get(context.Background(), &auth.SessionValue{Value: sessVal})
	if err != nil {
		return nil, errors.GRPCError(err)
	}

	return models.GrpcSessionToModel(sess), nil
}

func (ac *AuthClient) Delete(sessVal string) *errors.Error {
	_, err := ac.client.Delete(context.Background(), &auth.SessionValue{Value: sessVal})
	if err != nil {
		return errors.GRPCError(err)
	}
	return nil
}

func (ac *AuthClient) Check(sessValue string) (*models.Session, *errors.Error) {
	sess, err := ac.client.Check(context.Background(), &auth.SessionValue{Value: sessValue})
	if err != nil {
		return nil, errors.GRPCError(err)
	}
	return models.GrpcSessionToModel(sess), nil
}
