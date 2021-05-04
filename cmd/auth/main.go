package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2021_1_YSNP/configs"
	authGRPC "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/auth/delivery/grpc"
	authRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/auth/repository/tarantool"
	authUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/auth/usecase"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/databases"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/interceptor"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/auth"
)

func main() {
	configs, err := configs.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	tarantoolDB, err := databases.NewTarantool(configs.GetTarantoolUser(), configs.GetTarantoolPassword(), configs.GetTarantoolConfig())
	if err != nil {
		log.Fatal(err)
	}

	ar := authRepo.NewSessionRepository(tarantoolDB.GetDatabase())
	au := authUsecase.NewSessionUsecase(ar)
	handler := authGRPC.NewAuthHandlerServer(au)

	lis, err := net.Listen("tcp", fmt.Sprint(configs.GetAuthHost(), ":", configs.GetAuthPort()))
	if err != nil {
		log.Fatalln("Can't listen session microservice port", err)
	}
	defer lis.Close()

	logger := logger.NewLogger(configs.GetLoggerMode())
	logger.StartServerLog(configs.GetAuthHost(), configs.GetAuthPort())
	ic := interceptor.NewInterceptor(logger.GetLogger())

	server := grpc.NewServer(
		grpc.UnaryInterceptor(ic.ServerLogInterceptor),
	)
	auth.RegisterAuthHandlerServer(server, handler)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
