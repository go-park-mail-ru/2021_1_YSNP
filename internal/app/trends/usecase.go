package trends

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

type TrendsUsecase interface {
	InsertOrUpdate(ui *models.UserInterested) *errors.Error
}
