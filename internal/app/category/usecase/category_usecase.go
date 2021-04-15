package usecase

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

type CategoryUsecase struct {
	categoryRepo category.CategoryRepository
}

func NewCategoryUsecase(repo category.CategoryRepository) category.CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo: repo,
	}
}

func (cat *CategoryUsecase) GetAllCategories() ([]*models.Category, *errors.Error) {
	categories, err := cat.categoryRepo.GetAllCategories()
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return categories, nil
}
