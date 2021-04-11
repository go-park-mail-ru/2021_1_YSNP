package search

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

//go:generate mockgen -destination=./mocks/mock_search_repo.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search SearchRepository

type SearchRepository interface {
	SelectByFilter(data *models.Search) ([]*models.ProductListData, error)
}
