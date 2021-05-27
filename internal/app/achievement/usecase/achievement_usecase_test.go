package usecase

import (
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/achievement/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

func TestAchievementUsecase_GetUserAchievements(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	achRepo := mock.NewMockAchievementRepository(ctrl)
	achUcase := NewAchievementUsecase(achRepo)

	achRepo.EXPECT().GetUserAchievements(gomock.Eq(0)).Return([]*models.Achievement{}, nil)

	_, err := achUcase.GetUserAchievements(0)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestAchievementUsecase_GetUserAchievements_Error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	achRepo := mock.NewMockAchievementRepository(ctrl)
	achUcase := NewAchievementUsecase(achRepo)

	achRepo.EXPECT().GetUserAchievements(gomock.Eq(0)).Return(nil, sql.ErrConnDone)

	_, err := achUcase.GetUserAchievements(0)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}
