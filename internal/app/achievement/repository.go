package achievement

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

//go:generate mockgen -destination=./mocks/mock_achive_repo.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/achievement AchievementRepository

type AchievementRepository interface {
	GetUserAchievements(userId int, loggedUser int) ([]*models.Achievement, error)
}
