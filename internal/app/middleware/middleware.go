package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	LogrusLogger *logrus.Entry
	sessUcase    session.SessionUsecase
	userUcase    user.UserUsecase
}

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	ContextUserID = contextKey("userID")
	ContextLogger = contextKey("logger")
)

func NewMiddleware(sessUcase session.SessionUsecase, userUcase user.UserUsecase) *Middleware {
	return &Middleware{
		sessUcase: sessUcase,
		userUcase: userUcase,
	}
}

func CorsControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "*")

		switch req.Header.Get("Origin") {
		case "http://localhost:3000":
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		case "https://ykoya.ru":
			w.Header().Set("Access-Control-Allow-Origin", "https://ykoya.ru")
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if req.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, req)
	})
}

func (m *Middleware) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.LogrusLogger = m.LogrusLogger.WithFields(logrus.Fields{
			"method":  r.Method,
			"path":    r.URL.Path,
			"work_id": uuid.New(),
		})
		m.LogrusLogger.Info("Get connection")

		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextLogger, m.LogrusLogger)
		start := time.Now()
		next.ServeHTTP(w, r.WithContext(ctx))

		m.LogrusLogger.WithFields(logrus.Fields{
			"work_time": time.Since(start),
		}).Info("Fulfilled connection")
	})
}

func (m *Middleware) CheckAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		session, errE := m.sessUcase.Check(cookie.Value)
		if errE != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextUserID, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
