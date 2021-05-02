package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2021_1_YSNP/configs"
	authGRPC "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/auth/delivery/grpc"
	authRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/auth/repository/tarantool"
	authUcase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/auth/usecase"
	databases2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/databases"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/auth"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	configs, err := configs.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	tarantoolDB, err := databases2.NewTarantool(configs.GetTarantoolUser(), configs.GetTarantoolPassword(), configs.GetTarantoolConfig())
	if err != nil {
		log.Fatal(err)
	}

	ar := authRepo.NewSessionRepository(tarantoolDB.GetDatabase())
	au := authUcase.NewSessionUsecase(ar)
	handler := authGRPC.NewAuthHandlerServer(au)

	lis, err := net.Listen("tcp", fmt.Sprint(configs.GetAuthHost(), ":", configs.GetAuthPort()))

	if err != nil {
		log.Fatalln("Can't listen session microservice port", err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	auth.RegisterAuthHandlerServer(server, handler)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}