package search

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

type SearchRepository interface {
	SelectByFilter(data *models.Search)([]*models.ProductData, error)
}
