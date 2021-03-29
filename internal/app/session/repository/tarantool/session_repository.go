package repository

import (
	"encoding/json"
	"errors"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session"
	"github.com/tarantool/go-tarantool"
)

type SessionRepository struct {
	dbConn *tarantool.Connection
}

func NewSessionRepository(conn *tarantool.Connection) session.SessionRepository {
	return &SessionRepository{
		dbConn: conn,
	}
}

func (sr *SessionRepository) Insert(session *models.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	dataStr := string(data)

	//resp, err := sr.dbConn.Eval("return new_session(...)", []interface{}{session.Value, dataStr})
	_, err = sr.dbConn.Insert("sessions", []interface{}{session.Value, dataStr})
	if err != nil {
		return err
	}

	return nil
}

func (sr *SessionRepository) SelectByValue(sessValue string) (*models.Session, error) {
	resp, err := sr.dbConn.Call("check_session", []interface{}{sessValue})
	if err != nil {
		return nil, err
	}

	data := resp.Data[0]
	if data == nil {
		return &models.Session{}, nil
	}

	sessionDataSlice, ok := data.([]interface{})
	if !ok {
		return nil, errors.New("cannot cast data")
	}

	sessionData, ok := sessionDataSlice[1].(string)
	if !ok {
		return nil, errors.New("cannot cast")
	}

	sess := &models.Session{}
	err = json.Unmarshal([]byte(sessionData), sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (sr *SessionRepository) DeleteByValue(sessionValue string) error {
	// panic("implement me")
	return nil
}
