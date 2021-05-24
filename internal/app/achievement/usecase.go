package achievement

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

type AchievementUsecase interface {
	GetUserAchievements(userId int) ([]*models.Achievement, *errors.Error)
}
