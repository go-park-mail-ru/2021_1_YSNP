package product

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

type CategoryUsecase interface {
	GetCategory() ([]*models.Category, *errors.Error)
}
