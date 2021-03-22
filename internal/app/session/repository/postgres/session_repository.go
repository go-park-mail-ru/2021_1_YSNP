package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session"
)

type SessionRepository struct {
	dbConn *sql.DB
}

func NewSessionRepository(conn *sql.DB) session.SessionRepository {
	return &SessionRepository{
		dbConn: conn,
	}
}

func (sr *SessionRepository) Insert(session *models.Session) error {
	panic("implement me")
}

func (sr *SessionRepository) SelectByValue(sessValue string) (*models.Session, error) {
	panic("implement me")
}

func (sr *SessionRepository) DeleteByValue(sessionValue string) error {
	panic("implement me")
}