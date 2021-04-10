package product

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

type CategoryRepository interface {
	GetCategory() ([]*models.Category, error)
}
