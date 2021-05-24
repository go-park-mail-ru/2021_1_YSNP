package achievement

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

//go:generate mockgen -destination=./mocks/mock_product_repo.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product ProductRepository

type AchievementRepository interface {
	GetUserAchievements(userId int) ([]*models.Achievement, error)
}
