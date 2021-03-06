package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/metrics"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	log "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	responseObserver "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/utils"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
)

const CsrfKey = "ba6f7ee3-84d8-4f68-aaa5-5ef7c1823aa4"

type Middleware struct {
	logrusLogger *logrus.Entry
	sessUcase    session.SessionUsecase
	userUcase    user.UserUsecase
	metricsM     *metrics.Metrics
}

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	ContextUserID = contextKey("userID")
	ContextLogger = contextKey("logger")
)

func NewMiddleware(sessUcase session.SessionUsecase, userUcase user.UserUsecase, metrics *metrics.Metrics) *Middleware {
	return &Middleware{
		sessUcase: sessUcase,
		userUcase: userUcase,
		metricsM:  metrics,
	}
}

func (m *Middleware) NewLogger(logger *logrus.Entry) {
	m.logrusLogger = logger
}

func CorsControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "X-CSRF-Token, sentry-trace")

		switch req.Header.Get("Origin") {
		case "http://localhost:3000":
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		case "https://ykoya.ru":
			w.Header().Set("Access-Control-Allow-Origin", "https://ykoya.ru")
		}

		w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
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
		m.logrusLogger = m.logrusLogger.WithFields(logrus.Fields{
			"method":  r.Method,
			"path":    r.URL.Path,
			"work_id": uuid.New(),
		})
		m.logrusLogger.Info("Get connection")

		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextLogger, m.logrusLogger)
		start := time.Now()

		o := &responseObserver.ResponseObserver{ResponseWriter: w}
		next.ServeHTTP(o, r.WithContext(ctx))

		m.logrusLogger.WithFields(logrus.Fields{
			"work_time": time.Since(start),
		}).Info("Fulfilled connection")

		if r.URL.Path != "/metrics" {
			currentRoute := mux.CurrentRoute(r)
			pathTemplate, err := currentRoute.GetPathTemplate()
			if err != nil {
				m.logrusLogger.Warn("current path template")
				return
			}

			m.metricsM.Hits.WithLabelValues(strconv.Itoa(o.Status), pathTemplate, r.Method).Inc()
			m.metricsM.Timings.WithLabelValues(
				strconv.Itoa(o.Status), pathTemplate, r.Method).Observe(float64(start.Second()))
		}

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

func (m *Middleware) SetCSRFToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		next.ServeHTTP(w, r)
	}
}

func (m *Middleware) CSFRErrorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, ok := r.Context().Value(ContextLogger).(*logrus.Entry)
		if !ok {
			logger = log.GetDefaultLogger()
			logger.Warn("no logger")
		}

		errE := errors.Cause(errors.InvalidCSRFToken)
		logger.Error(errE.Message)

		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
	}
}
