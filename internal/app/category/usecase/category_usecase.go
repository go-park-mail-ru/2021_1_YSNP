package usecase

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category"
)

type CategoryUsecase struct {
	categoryRepo category.CategoryRepository
}

func NewProductUsecase(repo category.CategoryUsecase) category.CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo: repo,
	}
}

func(cat *CategoryUsecase) GetCategory() ([]*models.Category, *errors.Error) {
	return nil, nil
}