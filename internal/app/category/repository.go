package category

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

//go:generate mockgen -destination=./mocks/mock_category_repo.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category CategoryRepository

type CategoryRepository interface {
	SelectCategories() ([]*models.Category, error)
}
