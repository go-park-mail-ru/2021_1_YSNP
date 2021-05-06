package trends

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

//go:generate mockgen -destination=./mocks/mock_trends_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends TrendsUsecase

type TrendsUsecase interface {
	InsertOrUpdate(ui *models.UserInterested) *errors.Error
}
