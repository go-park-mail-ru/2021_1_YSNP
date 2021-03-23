package middleware

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Middleware struct {
	LogrusLogger *logrus.Entry
	sessUcase session.SessionUsecase
	userUcase user.UserUsecase
}

func NewMiddleware (sessUcase session.SessionUsecase, userUcase user.UserUsecase) *Middleware {
	return &Middleware{
		sessUcase:    sessUcase,
		userUcase:    userUcase,
	}
}


func CorsControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "*")

		switch req.Header.Get("Origin") {
		case "http://localhost:3000":
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		case "http://89.208.199.170:3000":
			w.Header().Set("Access-Control-Allow-Origin", "http://89.208.199.170:3000")
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
		start := time.Now()
		next.ServeHTTP(w, r)

		m.LogrusLogger.WithFields(logrus.Fields{
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start),
		}).Info(r.URL.Path)
	})
}

func (m *Middleware) CheckAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
		}

		session, err := m.sessUcase.Check(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
