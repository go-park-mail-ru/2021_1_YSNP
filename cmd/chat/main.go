package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2021_1_YSNP/configs"
	chatGRPC "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/chat/delivery/grpc"
	chatRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/chat/repository/postgres"
	chatUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/chat/usecase"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/databases"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/interceptor"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/chat"
)

func main() {
	configs, err := configs.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	postgresDB, err := databases.NewPostgres(configs.GetPostgresConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	cr := chatRepo.NewChatRepository(postgresDB.GetDatabase())
	cu := chatUsecase.NewChatUsecase(cr)
	handler := chatGRPC.NewChatServer(cu)

	lis, err := net.Listen("tcp", fmt.Sprint(configs.GetChatHost(), ":", configs.GetChatPort()))
	if err != nil {
		log.Fatalln("Can't listen chat microservice port", err)
	}
	defer lis.Close()

	logger := logger.NewLogger(configs.GetLoggerMode())
	logger.StartServerLog(configs.GetChatHost(), configs.GetChatPort())
	ic := interceptor.NewInterceptor(logger.GetLogger())

	server := grpc.NewServer(
		grpc.UnaryInterceptor(ic.ServerLogInterceptor),
	)
	chat.RegisterChatServer(server, handler)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
