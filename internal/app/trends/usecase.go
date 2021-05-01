package trends

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	errors2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

type TrendsUsecase interface {
	InsertOrUpdate(ui *models.UserInterested) *errors2.Error
}
