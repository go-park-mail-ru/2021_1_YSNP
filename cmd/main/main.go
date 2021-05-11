package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"google.golang.org/grpc"

	appConfig "github.com/go-park-mail-ru/2021_1_YSNP/configs"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/metrics"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/databases"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/interceptor"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
	_ "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/validator"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/websocket"

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

	trendsHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends/delivery/http"
	trendsRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends/repository/tarantool"
	trendsUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends/usecase"

	sessHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/delivery/http"
	sessUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/usecase"

	chatHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/chat/delivery/http"
	chatWSHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/chat/delivery/websocket"
	chatUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/chat/usecase"
)

func main() {
	configs, err := appConfig.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	appConfig.Configs = configs

	postgresDB, err := databases.NewPostgres(configs.GetPostgresConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	tarantoolDB, err := databases.NewTarantool(configs.GetTarantoolUser(), configs.GetTarantoolPassword(), configs.GetTarantoolConfig())
	if err != nil {
		log.Fatal(err)
	}

	trendsRepo := trendsRepo.NewTrendsRepository(tarantoolDB.GetDatabase(), postgresDB.GetDatabase())
	userRepo := userRepo.NewUserRepository(postgresDB.GetDatabase())
	prodRepo := productRepo.NewProductRepository(postgresDB.GetDatabase())
	searchRepo := searchRepo.NewSearchRepository(postgresDB.GetDatabase())
	categoryRepo := categoryRepo.NewCategoryRepository(postgresDB.GetDatabase())
	uploadRepo := uploadRepo.NewUploadRepository()

	userUcase := userUsecase.NewUserUsecase(userRepo, uploadRepo)
	prodUcase := productUsecase.NewProductUsecase(prodRepo, uploadRepo, trendsRepo)
	searchUcase := searchUsecase.NewSearchUsecase(searchRepo)
	categoryUsecase := categoryUsecase.NewCategoryUsecase(categoryRepo)
	trendsUsecase := trendsUsecase.NewTrendsUsecase(trendsRepo)

	logger := logger.NewLogger(configs.GetLoggerMode())
	logger.StartServerLog(configs.GetMainHost(), configs.GetMainPort())
	ic := interceptor.NewInterceptor(logger.GetLogger())

	jaeger, err := metrics.NewJaeger("client")
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}

	jaeger.SetGlobalTracer()
	defer jaeger.Close()

	sessionGRPCConn, err := grpc.Dial(
		fmt.Sprint(configs.GetAuthHost(), ":", configs.GetAuthPort()),
		grpc.WithChainUnaryInterceptor(traceutils.OpenTracingClientInterceptor(jaeger.GetTracer()), ic.ClientLogInterceptor),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer sessionGRPCConn.Close()
	sessUcase := sessUsecase.NewAuthClient(sessionGRPCConn)

	chatGRPCConn, err := grpc.Dial(
		fmt.Sprint(configs.GetChatHost(), ":", configs.GetChatPort()),
		grpc.WithChainUnaryInterceptor(traceutils.OpenTracingClientInterceptor(jaeger.GetTracer()), ic.ClientLogInterceptor),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer chatGRPCConn.Close()
	chatUcase := chatUsecase.NewChatClient(chatGRPCConn)

	userHandler := userHandler.NewUserHandler(userUcase, sessUcase)
	prodHandler := productHandler.NewProductHandler(prodUcase)
	searchHandler := searchHandler.NewSearchHandler(searchUcase)
	categoryHandler := categoryHandler.NewCategoryHandler(categoryUsecase)
	trendsHandler := trendsHandler.NewTrendsHandler(trendsUsecase)

	chatHandler := chatHandler.NewChatHandler(chatUcase)
	chatWSHandler := chatWSHandler.NewChatWSHandler(chatUcase)
	sessHandler := sessHandler.NewSessionHandler(sessUcase, userUcase)

	router := mux.NewRouter()
	metricsProm := metrics.NewMetrics(router)

	mw := middleware.NewMiddleware(sessUcase, userUcase, metricsProm)
	mw.NewLogger(logger.GetLogger())

	router.Use(middleware.CorsControlMiddleware)
	router.Use(mw.AccessLogMiddleware)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(csrf.Protect([]byte(middleware.CsrfKey),
		csrf.ErrorHandler(mw.CSFRErrorHandler())))

	wsSrv := websocket.NewWSServer(logger)
	wsSrv.Run()
	defer wsSrv.Stop()

	userHandler.Configure(api, mw)
	sessHandler.Configure(api, mw)
	prodHandler.Configure(api, router, mw)
	trendsHandler.Configure(api, mw)
	searchHandler.Configure(api, mw)
	categoryHandler.Configure(api, mw)
	chatHandler.Configure(api, mw, wsSrv)
	chatWSHandler.Configure(api, mw, wsSrv)

	server := http.Server{
		Addr:         fmt.Sprint(":", configs.GetMainPort()),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
