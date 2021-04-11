package search

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

//go:generate mockgen -destination=./mocks/mock_search_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search SearchUsecase

type SearchUsecase interface {
	SelectByFilter(data *models.Search) ([]*models.ProductListData, *errors.Error)
}
