package grpc

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/auth"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	proto "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/auth"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type AuthHandlerServer struct {
	authUcase auth.SessionUsecase
	logrusLogger *logrus.Entry
}

func NewAuthHandlerServer(au auth.SessionUsecase) *AuthHandlerServer {
	return &AuthHandlerServer{
		authUcase: au,
	}
}

func (a *AuthHandlerServer) NewLogger(logger *logrus.Entry) {
	a.logrusLogger = logger
}

func (a *AuthHandlerServer) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	a.logrusLogger = a.logrusLogger.WithFields(logrus.Fields{
		"method":  info.FullMethod,
		"request": req,
		"work_id": uuid.New(),
		"metadata": md,
	})
	a.logrusLogger.Info("Get connection")

	start := time.Now()


	reply, err := handler(ctx, req)

	a.logrusLogger.WithFields(logrus.Fields{
		"work_time": time.Since(start),
		"reply": reply,
		"error": err,
	}).Info("Fulfilled connection")
	return reply, err
}

func (a *AuthHandlerServer) Create (ctx context.Context, sess *proto.Session) (*emptypb.Empty, error) {
	if err := a.authUcase.Create(models.GrpcSessionToModel(sess)); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Code(err.ErrorCode), err.Message)
	}
	return &emptypb.Empty{}, nil
}

func (a *AuthHandlerServer) Get(ctx context.Context, sessVal *proto.SessionValue) (*proto.Session, error) {
	sess, err := a.authUcase.Get(sessVal.GetValue())
	if err != nil {
		return nil, status.Error(codes.Code(err.ErrorCode), err.Message)
	}
	return models.ModelSessionToGrpc(sess), nil
}

func (a *AuthHandlerServer) Check(ctx context.Context, sessVal *proto.SessionValue) (*proto.Session, error) {
	sess, err := a.authUcase.Check(sessVal.GetValue())
	if err != nil {
		return nil, status.Error(codes.Code(err.ErrorCode), err.Message)
	}
	return models.ModelSessionToGrpc(sess), nil
}

func (a *AuthHandlerServer) Delete(ctx context.Context, sessVal *proto.SessionValue) (*emptypb.Empty, error) {
	if err := a.authUcase.Delete(sessVal.GetValue()); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Code(err.ErrorCode), err.Message)
	}
	return &emptypb.Empty{}, nil
}