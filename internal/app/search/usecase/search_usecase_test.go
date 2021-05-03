package usecase

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search/mocks"
	errors "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearchUsecase_SelectByFilter_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	searchRepo := mock.NewMockSearchRepository(ctrl)
	searchUcase := NewSearchUsecase(searchRepo)

	search := &models.Search{
		Category:   "Товар",
		Date:       "",
		FromAmount: 0,
		ToAmount:   0,
		Radius:     0,
		Latitude:   0,
		Longitude:  0,
		Search:     "",
		Sorting:    "",
		From:       0,
		Count:      0,
	}

	var userID uint64 = 0

	prod := &models.ProductListData{
		ID:         0,
		Name:       "Товар",
		Date:       "2100-02-02",
		Amount:     1000,
		LinkImages: nil,
		UserLiked:  false,
		Tariff:     0,
	}

	searchRepo.EXPECT().SelectByFilter(&userID, search).Return([]*models.ProductListData{prod}, nil)

	res, err := searchUcase.SelectByFilter(&userID, search)
	assert.Equal(t, res[0], prod)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	searchRepo.EXPECT().SelectByFilter(&userID, search).Return( nil, sql.ErrConnDone)

	_, err = searchUcase.SelectByFilter(&userID, search)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestSearchUsecase_SelectByFilter_EmptySearch(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	searchRepo := mock.NewMockSearchRepository(ctrl)
	searchUcase := NewSearchUsecase(searchRepo)

	search := &models.Search{
		Category:   "Товар",
		Date:       "",
		FromAmount: 0,
		ToAmount:   0,
		Radius:     0,
		Latitude:   0,
		Longitude:  0,
		Search:     "",
		Sorting:    "",
		From:       0,
		Count:      0,
	}

	var userID uint64 = 0

	searchRepo.EXPECT().SelectByFilter(&userID, search).Return([]*models.ProductListData{}, nil)

	_, err := searchUcase.SelectByFilter(&userID, search)
	assert.Equal(t, err, errors.Cause(errors.EmptySearch))
}
