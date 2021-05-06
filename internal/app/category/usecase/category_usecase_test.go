package usecase

import (
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

func TestCategoryUsecase_GetAllCategories(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	catRepo := mock.NewMockCategoryRepository(ctrl)
	catUcase := NewCategoryUsecase(catRepo)

	catRepo.EXPECT().SelectCategories().Return([]*models.Category{}, nil)

	_, err := catUcase.GetAllCategories()
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestCategoryUsecase_GetAllCategories_Error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	catRepo := mock.NewMockCategoryRepository(ctrl)
	catUcase := NewCategoryUsecase(catRepo)

	catRepo.EXPECT().SelectCategories().Return(nil, sql.ErrConnDone)

	_, err := catUcase.GetAllCategories()
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}
