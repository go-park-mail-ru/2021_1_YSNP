package usecase

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/achievement"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

type AchievementUsecase struct {
	achRepo achievement.AchievementRepository
}

func NewAchievementUsecase(repo achievement.AchievementRepository) achievement.AchievementUsecase {
	return &AchievementUsecase{
		achRepo: repo,
	}
}

func (au *AchievementUsecase) GetUserAchievements(userId int) ([]*models.Achievement, *errors.Error){
	achievement, err := au.achRepo.GetUserAchievements(userId)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return achievement, nil
}
