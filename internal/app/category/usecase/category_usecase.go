package usecase

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	errors2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

type CategoryUsecase struct {
	categoryRepo category.CategoryRepository
}

func NewCategoryUsecase(repo category.CategoryRepository) category.CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo: repo,
	}
}

func (cat *CategoryUsecase) GetAllCategories() ([]*models.Category, *errors2.Error) {
	categories, err := cat.categoryRepo.SelectCategories()
	if err != nil {
		return nil, errors2.UnexpectedInternal(err)
	}

	return categories, nil
}
