package usecase

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category"
)

type CategoryUsecase struct {
	categoryRepo category.CategoryRepository
}

func NewCategoryUsecase(repo category.CategoryRepository) category.CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo: repo,
	}
}

func(cat *CategoryUsecase) GetCategory() ([]*models.Category, *errors.Error) {
	categories, err := cat.categoryRepo.GetCategory()
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}
	
	return categories, nil
}