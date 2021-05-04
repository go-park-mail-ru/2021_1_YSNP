package interceptor

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Interceptor struct {
	logrusLogger *logrus.Entry
}

func NewInterceptor(logger *logrus.Entry) *Interceptor {
	return &Interceptor{
		logrusLogger: logger,
	}
}

func (i *Interceptor) ServerLogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	i.logrusLogger = i.logrusLogger.WithFields(logrus.Fields{
		"method":  info.FullMethod,
		"request": req,
		"work_id": uuid.New(),
	})
	i.logrusLogger.Info("Get connection")
	start := time.Now()

	reply, err := handler(ctx, req)
	if err != nil {
		i.logrusLogger.WithFields(logrus.Fields{
			"req":   req,
			"error": err.Error(),
		}).Error("Error occurred")
	}

	i.logrusLogger.WithFields(logrus.Fields{
		"work_time": time.Since(start),
	}).Info("Fulfilled connection")

	return reply, err
}

func (i *Interceptor) ClientLogInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	i.logrusLogger = i.logrusLogger.WithFields(logrus.Fields{
		"method":  method,
		"request": req,
		"work_id": uuid.New(),
	})
	i.logrusLogger.Info("Connect to ", cc.Target())
	start := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		i.logrusLogger.WithFields(logrus.Fields{
			"client_conn": cc.Target(),
			"error":       err.Error(),
		}).Error("Error occurred")
	}

	i.logrusLogger.WithFields(logrus.Fields{
		"work_time": time.Since(start),
	}).Info("Fulfilled connection")

	return err
}
