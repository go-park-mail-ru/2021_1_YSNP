package main

import (
	"fmt"
	sessHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/delivery/http"
	sessUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/usecase"
	interceptor2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/interceptor"
	middleware2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"

	chatHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/chat/delivery/http"
	chatWSHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/chat/delivery/websocket"
	chatUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/chat/usecase"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/websocket"

	databases2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/databases"
	logger2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2021_1_YSNP/configs"
	_ "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/validator"

	categoryHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category/delivery/http"
	categoryRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category/repository/postgres"
	categoryUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category/usecase"

	userHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/delivery/http"
	userRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/repository/postgres"
	userUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/usecase"

	productHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/delivery/http"
	productRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/repository/postgres"
	productUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/usecase"
	searchHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search/delivery/http"
	searchRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search/repository/postgres"
	searchUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search/usecase"

	uploadRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/upload/repository/FileSystem"
)

func main() {
	configs, err := configs.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	postgresDB, err := databases2.NewPostgres(configs.GetPostgresConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	userRepo := userRepo.NewUserRepository(postgresDB.GetDatabase())
	prodRepo := productRepo.NewProductRepository(postgresDB.GetDatabase())
	searchRepo := searchRepo.NewSearchRepository(postgresDB.GetDatabase())
	categoryRepo := categoryRepo.NewCategoryRepository(postgresDB.GetDatabase())
	uploadRepo := uploadRepo.NewUploadRepository()

	userUcase := userUsecase.NewUserUsecase(userRepo, uploadRepo)
	prodUcase := productUsecase.NewProductUsecase(prodRepo, uploadRepo)
	searchUcase := searchUsecase.NewSearchUsecase(searchRepo)
	categoryUsecase := categoryUsecase.NewCategoryUsecase(categoryRepo)

	logger := logger2.NewLogger(configs.GetLoggerMode())
	logger.StartServerLog(configs.GetServerHost(), configs.GetServerPort())

	ic := interceptor2.NewInterceptor(logger.GetLogger())

	sessionGRPCConn, err := grpc.Dial(
		fmt.Sprint(configs.GetAuthHost(), ":", configs.GetAuthPort()),
		grpc.WithUnaryInterceptor(ic.ClientLogInterceptor),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer sessionGRPCConn.Close()
	sessUcase := sessUsecase.NewAuthClient(sessionGRPCConn)

	chatGRPCConn, err := grpc.Dial(
		fmt.Sprint(configs.GetChatHost(), ":", configs.GetChatPort()),
		grpc.WithUnaryInterceptor(ic.ClientLogInterceptor),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer chatGRPCConn.Close()
	chatUcase := chatUsecase.NewChatClient(chatGRPCConn)

	userHandler := userHandler.NewUserHandler(userUcase, sessUcase)
	sessHandler := sessHandler.NewSessionHandler(sessUcase, userUcase)
	prodHandler := productHandler.NewProductHandler(prodUcase)
	searchHandler := searchHandler.NewSearchHandler(searchUcase)
	categoryHandler := categoryHandler.NewCategoryHandler(categoryUsecase)
	chatHandler := chatHandler.NewChatHandler(chatUcase)
	chatWSHandler := chatWSHandler.NewChatWSHandler(chatUcase)

	router := mux.NewRouter()

	mw := middleware2.NewMiddleware(sessUcase, userUcase)
	mw.NewLogger(logger.GetLogger())

	router.Use(middleware2.CorsControlMiddleware)
	router.Use(mw.AccessLogMiddleware)

	api := router.PathPrefix("/api/v1").Subrouter()
	//api.Use(csrf.Protect([]byte(middleware.CsrfKey),
	//	csrf.ErrorHandler(mw.CSFRErrorHandler())))

	wsSrv := websocket.NewWSServer()
	wsSrv.Run()
	defer wsSrv.Stop()

	userHandler.Configure(api, mw)
	sessHandler.Configure(api, mw)
	prodHandler.Configure(api, router, mw)
	searchHandler.Configure(api, mw)
	categoryHandler.Configure(api, mw)
	chatHandler.Configure(api, mw, wsSrv)
	chatWSHandler.Configure(api, mw, wsSrv)

	server := http.Server{
		Addr:         fmt.Sprint(":", configs.GetServerPort()),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
