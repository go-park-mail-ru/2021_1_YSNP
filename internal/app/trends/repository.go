package trends

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

//go:generate mockgen -destination=./mocks/mock_trends_repo.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends TrendsRepository

type TrendsRepository interface {
	InsertOrUpdate(ui *models.Trends) error
	CreateTrendsProducts(userID uint64) error
	GetTrendsProducts(userID uint64) ([]uint64, error)
	GetRecommendationProducts(productID uint64, userID uint64) ([]uint64, error)
}

