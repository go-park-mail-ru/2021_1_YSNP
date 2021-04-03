package search

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

type SearchUsecase interface {
	SelectByFilter(data *models.Search) ([]*models.ProductListData, *errors.Error)
}
