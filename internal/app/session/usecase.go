package session

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	errors "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

//go:generate mockgen -destination=./mocks/mock_session_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session SessionUsecase
//go:generate mockgen -destination=./mocks/mock_session_client.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/auth AuthHandlerClient

type SessionUsecase interface {
	Create(sess *models.Session) *errors.Error
	Get(sessValue string) (*models.Session, *errors.Error)
	Delete(sessionValue string) *errors.Error
	Check(sessValue string) (*models.Session, *errors.Error)
}
