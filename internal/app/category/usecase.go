package category

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

//go:generate mockgen -destination=./mocks/mock_category_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category CategoryUsecase

type CategoryUsecase interface {
	GetAllCategories() ([]*models.Category, *errors.Error)
}
