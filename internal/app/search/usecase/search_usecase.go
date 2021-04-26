package usecase

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search"
	errors2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

type SearchUsecase struct {
	searchRepo search.SearchRepository
}

func NewSearchUsecase(repo search.SearchRepository) search.SearchUsecase {
	return &SearchUsecase{
		searchRepo: repo,
	}
}

func (su *SearchUsecase) SelectByFilter(userID *uint64, data *models.Search) ([]*models.ProductListData, *errors2.Error) {
	res, err := su.searchRepo.SelectByFilter(userID, data)

	if err != nil {
		return nil, errors2.UnexpectedInternal(err)
	}

	if len(res) == 0 {
		return nil, errors2.Cause(errors2.EmptySearch)
	}
	return res, nil
}
