package achievement

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

//go:generate mockgen -destination=./mocks/mock_achive_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/achievement AchievementUsecase

type AchievementUsecase interface {
	GetUserAchievements(userId int) ([]*models.Achievement, *errors.Error)
}
