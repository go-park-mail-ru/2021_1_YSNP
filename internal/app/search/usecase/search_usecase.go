package usecase

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search"
)

type SearchUsecase struct {
	searchRepo search.SearchRepository
}

func NewSessionUsecase(repo search.SearchRepository) search.SearchUsecase {
	return &SearchUsecase {
		searchRepo: repo,
	}
}

func (su *SearchUsecase) SelectByFilter(data *models.Search) ([]*models.ProductData, *errors.Error) {
	res, _ := su.searchRepo.SelectByFilter(data)
	return res, nil
}
