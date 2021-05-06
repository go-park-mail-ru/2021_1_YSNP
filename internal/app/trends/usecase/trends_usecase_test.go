package usecase

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrendsUsecase_InsertOrUpdate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trendsRepo := mock.NewMockTrendsRepository(ctrl)
	trendsUcase := NewTrendsUsecase(trendsRepo)

	trendsRepo.EXPECT().InsertOrUpdate(gomock.Any()).Return(sql.ErrConnDone)

	err := trendsUcase.InsertOrUpdate(&models.UserInterested{
		UserID: 1,
		Text:   "Словов",
	})
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}