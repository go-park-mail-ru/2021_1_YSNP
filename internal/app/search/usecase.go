package search

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

//go:generate mockgen -destination=./mocks/mock_search_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search SearchUsecase

type SearchUsecase interface {
	SelectByFilter(userID *uint64, data *models.Search) ([]*models.ProductListData, *errors.Error)
}
