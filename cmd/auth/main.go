package main

import (
	"fmt"
	"log"
	"net"

	traceutils "github.com/opentracing-contrib/go-grpc"
	"google.golang.org/grpc"

	authGRPC "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/auth/delivery/grpc"
	authRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/auth/repository/tarantool"
	authUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/auth/usecase"

	"github.com/go-park-mail-ru/2021_1_YSNP/configs"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/metrics"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/databases"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/interceptor"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/auth"
)

func main() {
	err := configs.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	tarantoolDB, err := databases.NewTarantool(configs.Configs.GetTarantoolUser(), configs.Configs.GetTarantoolPassword(), configs.Configs.GetTarantoolConfig())
	if err != nil {
		log.Fatal(err)
	}

	ar := authRepo.NewSessionRepository(tarantoolDB.GetDatabase())
	au := authUsecase.NewSessionUsecase(ar)
	handler := authGRPC.NewAuthHandlerServer(au)

	lis, err := net.Listen("tcp", fmt.Sprint(configs.Configs.GetAuthHost(), ":", configs.Configs.GetAuthPort()))
	if err != nil {
		log.Fatalln("Can't listen session microservice port", err)
	}
	defer lis.Close()

	logger := logger.NewLogger(configs.Configs.GetLoggerMode())
	logger.StartServerLog(configs.Configs.GetAuthHost(), configs.Configs.GetAuthPort())
	ic := interceptor.NewInterceptor(logger.GetLogger())

	jaeger, err := metrics.NewJaeger("auth")
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}

	jaeger.SetGlobalTracer()
	defer jaeger.Close()

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(traceutils.OpenTracingServerInterceptor(jaeger.GetTracer()), ic.ServerLogInterceptor),
	)
	auth.RegisterAuthHandlerServer(server, handler)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
