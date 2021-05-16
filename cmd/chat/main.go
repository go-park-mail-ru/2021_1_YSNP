package main

import (
	"fmt"
	"log"
	"net"

	traceutils "github.com/opentracing-contrib/go-grpc"
	"google.golang.org/grpc"

	chatGRPC "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/chat/delivery/grpc"
	chatRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/chat/repository/postgres"
	chatUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/chat/usecase"

	"github.com/go-park-mail-ru/2021_1_YSNP/configs"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/metrics"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/databases"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/interceptor"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/chat"
)

func main() {
	err := configs.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	postgresDB, err := databases.NewPostgres(configs.Configs.GetPostgresConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	cr := chatRepo.NewChatRepository(postgresDB.GetDatabase())
	cu := chatUsecase.NewChatUsecase(cr)
	handler := chatGRPC.NewChatServer(cu)

	lis, err := net.Listen("tcp", fmt.Sprint(configs.Configs.GetChatHost(), ":", configs.Configs.GetChatPort()))
	if err != nil {
		log.Fatalln("Can't listen chat microservice port", err)
	}
	defer lis.Close()

	logger := logger.NewLogger(configs.Configs.GetLoggerMode(), configs.Configs.GetChatPort())
	logger.StartServerLog(configs.Configs.GetChatHost(), configs.Configs.GetChatPort())
	ic := interceptor.NewInterceptor(logger.GetLogger())

	jaeger, err := metrics.NewJaeger("chat")
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}

	jaeger.SetGlobalTracer()
	defer jaeger.Close()

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(traceutils.OpenTracingServerInterceptor(jaeger.GetTracer()), ic.ServerLogInterceptor),
	)
	chat.RegisterChatServer(server, handler)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
