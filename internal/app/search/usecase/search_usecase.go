package usecase

import (
	"fmt"

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

func (su *SearchUsecase) SelectByFilter(data *models.Search) ([]*models.ProductListData, *errors.Error) {
	res, err := su.searchRepo.SelectByFilter(data)

	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}
	return res, nil
}
