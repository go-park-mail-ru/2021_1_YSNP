package usecase

import (
	"database/sql"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	errors "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
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
