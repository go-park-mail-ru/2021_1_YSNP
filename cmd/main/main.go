package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2021_1_YSNP/configs"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/databases"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	_ "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/validator"

	userHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/delivery/http"
	userRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/repository/postgres"
	userUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/usecase"

	sessionHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/delivery/http"
	sessionRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/repository/tarantool"
	sessionUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/usecase"

	productHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/delivery/http"
	productRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/repository/postgres"
	productUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/usecase"
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

	tarantoolDB, err := databases.NewTarantool(configs.GetTarantoolUser(), configs.GetTarantoolPassword(), configs.GetTarantoolConfig())
	if err != nil {
		log.Fatal(err)
	}

	userRepo := userRepo.NewUserRepository(postgresDB.GetDatabase())
	sessRepo := sessionRepo.NewSessionRepository(tarantoolDB.GetDatabase())
	prodRepo := productRepo.NewProductRepository(postgresDB.GetDatabase())

	userUcase := userUsecase.NewUserUsecase(userRepo)
	sessUcase := sessionUsecase.NewSessionUsecase(sessRepo)
	prodUcase := productUsecase.NewProductUsecase(prodRepo)

	userHandler := userHandler.NewUserHandler(userUcase, sessUcase)
	sessHandler := sessionHandler.NewSessionHandler(sessUcase, userUcase)
	prodHandler := productHandler.NewProductHandler(prodUcase)

	logger := logger.NewLogger(configs.GetLoggerMode())
	logger.StartServerLog(configs.GetServerHost(), configs.GetServerPort())

	router := mux.NewRouter()

	mw := middleware.NewMiddleware(sessUcase, userUcase, logger.GetLogger())
	router.Use(mw.AccessLogMiddleware)
	router.Use(middleware.CorsControlMiddleware)

	api := router.PathPrefix("/api/v1").Subrouter()

	userHandler.Configure(api, mw)
	sessHandler.Configure(api, mw)
	prodHandler.Configure(api, mw)

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
