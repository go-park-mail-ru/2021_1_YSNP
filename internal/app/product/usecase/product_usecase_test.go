package usecase

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/mocks"
	uMock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/upload/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"
)

var prodTest = &models.ProductData{
	ID:          0,
	Name:        "tovar",
	Date:        "",
	Amount:      10000,
	LinkImages:  []string{},
	Description: "Description product aaaaa",
	Category:    "0",
	OwnerID:     0,
}

func TestProductUsecase_Create_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	prodRepo.EXPECT().Insert(gomock.Eq(prodTest)).Return(nil)

	err := prodUcase.Create(prodTest)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestProductUsecase_GetByID_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	prodRepo.EXPECT().SelectByID(gomock.Eq(prodTest.ID)).Return(prodTest, nil)

	product, err := prodUcase.GetByID(prodTest.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, product, prodTest)
}

func TestProductUsecase_GetByID_ProductNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	prodRepo.EXPECT().SelectByID(gomock.Eq(prodTest.ID)).Return(nil, sql.ErrNoRows)

	_, err := prodUcase.GetByID(prodTest.ID)
	assert.Equal(t, err, errors.Cause(errors.ProductNotExist))
}

func TestProductUsecase_ListLatest_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	page := &models.Page{
		From:  0,
		Count: 20,
	}

	prodList := &models.ProductListData{
		ID:         0,
		Name:       "Product Name",
		Date:       "2013-3-3",
		Amount:     12000,
		LinkImages: nil,
		UserLiked:  false,
		Tariff:     0,
	}

	var userID uint64 = 0

	prodRepo.EXPECT().SelectLatest(&userID, gomock.Eq(page)).Return([]*models.ProductListData{prodList}, nil)

	list, err := prodUcase.ListLatest(&userID, page)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, list[0], prodList)
}

func TestProductUsecase_ListLatest_NoProduct(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	page := &models.Page{
		From:  0,
		Count: 20,
	}

	var userID uint64 = 0

	prodRepo.EXPECT().SelectLatest(&userID, gomock.Eq(page)).Return([]*models.ProductListData{}, nil)

	list, err := prodUcase.ListLatest(&userID, page)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, len(list), 0)
}

func TestProductUsecase_UserAdList_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	page := &models.Page{
		From:  0,
		Count: 20,
	}

	prodList := &models.ProductListData{
		ID:         0,
		Name:       "Product Name",
		Date:       "2013-3-3",
		Amount:     12000,
		LinkImages: nil,
		UserLiked:  false,
		Tariff:     0,
	}

	prodRepo.EXPECT().SelectUserAd(uint64(0), gomock.Eq(page)).Return([]*models.ProductListData{prodList}, nil)

	list, err := prodUcase.UserAdList(0, page)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, list[0], prodList)
}

func TestProductUsecase_UserAdList_NoProduct(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	page := &models.Page{
		From:  0,
		Count: 20,
	}

	prodRepo.EXPECT().SelectUserAd(uint64(0), gomock.Eq(page)).Return([]*models.ProductListData{}, nil)

	list, err := prodUcase.UserAdList(0, page)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, len(list), 0)
}

func TestProductUsecase_GetUserFavorite_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	page := &models.Page{
		From:  0,
		Count: 20,
	}

	prodList := &models.ProductListData{
		ID:         0,
		Name:       "Product Name",
		Date:       "2013-3-3",
		Amount:     12000,
		LinkImages: nil,
		UserLiked:  false,
		Tariff:     0,
	}

	prodRepo.EXPECT().SelectUserFavorite(uint64(0), gomock.Eq(page)).Return([]*models.ProductListData{prodList}, nil)

	list, err := prodUcase.GetUserFavorite(0, page)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, list[0], prodList)
}

func TestProductUsecase_GetUserFavorite_NoProduct(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	page := &models.Page{
		From:  0,
		Count: 20,
	}

	prodRepo.EXPECT().SelectUserFavorite(uint64(0), gomock.Eq(page)).Return([]*models.ProductListData{}, nil)

	list, err := prodUcase.GetUserFavorite(0, page)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, len(list), 0)
}

func TestProductUsecase_LikeProduct_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	prodRepo.EXPECT().InsertProductLike(uint64(0), uint64(0)).Return(nil)

	err := prodUcase.LikeProduct(0, 0)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestProductUsecase_LikeProduct(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	prodRepo.EXPECT().InsertProductLike(uint64(0), uint64(0)).Return(pgx.PgError{Code: "23505"})

	err := prodUcase.LikeProduct(0, 0)
	assert.Equal(t, err, errors.Cause(errors.ProductAlreadyLiked))
}

func TestProductUsecase_DislikeProduct_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	prodRepo.EXPECT().DeleteProductLike(uint64(0), uint64(0)).Return(nil)

	err := prodUcase.DislikeProduct(0, 0)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestProductUsecase_SetTariff_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	prodRepo.EXPECT().UpdateTariff(uint64(0), 0).Return(nil)

	err := prodUcase.SetTariff(0, 0)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestProductUsecase_UpdatePhoto_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(prodTest, nil)
	uploadRepo.EXPECT().InsertPhotos(gomock.Eq([]*multipart.FileHeader{}),"static/product/").Return([]string{}, nil)
	prodRepo.EXPECT().InsertPhoto(gomock.Any()).Return(nil)
	uploadRepo.EXPECT().RemovePhotos(gomock.Any()).Return(nil)

	prod, err := prodUcase.UpdatePhoto(prodTest.ID, uint64(0), []*multipart.FileHeader{})
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, prod, prodTest)
}

func TestProductUsecase_UpdatePhoto_NoProduct(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(nil, sql.ErrNoRows)

	_, err := prodUcase.UpdatePhoto(prodTest.ID, uint64(0), []*multipart.FileHeader{})
	assert.Equal(t, err, errors.Cause(errors.ProductNotExist))
}

func TestProductUsecase_UpdatePhoto_Error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo)

	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(prodTest, nil)
	uploadRepo.EXPECT().InsertPhotos(gomock.Any(), "static/product/").Return([]string{}, sql.ErrConnDone)

	_, err := prodUcase.UpdatePhoto(prodTest.ID, uint64(0), []*multipart.FileHeader{})
	assert.Equal(t, err.ErrorCode, errors.InternalError)
}
