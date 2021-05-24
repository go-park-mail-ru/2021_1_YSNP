package usecase

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	tMock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends/mocks"
	uMock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/upload/mocks"
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
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	prodRepo.EXPECT().Insert(gomock.Eq(prodTest)).Return(nil)

	err := prodUcase.Create(prodTest)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	prodRepo.EXPECT().Insert(gomock.Eq(prodTest)).Return(sql.ErrConnDone)

	err = prodUcase.Create(prodTest)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_GetByID_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	prodRepo.EXPECT().SelectByID(gomock.Eq(prodTest.ID)).Return(prodTest, nil)

	product, err := prodUcase.GetByID(prodTest.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, product, prodTest)

	//error
	prodRepo.EXPECT().SelectByID(gomock.Eq(prodTest.ID)).Return(nil, sql.ErrConnDone)

	product, err = prodUcase.GetByID(prodTest.ID)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_GetByID_ProductNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

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
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

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

	//error
	prodRepo.EXPECT().SelectLatest(&userID, gomock.Eq(page)).Return(nil, sql.ErrConnDone)

	_, err = prodUcase.ListLatest(&userID, page)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_ListLatest_NoProduct(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

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
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

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

	//error
	prodRepo.EXPECT().SelectUserAd(uint64(0), gomock.Eq(page)).Return(nil, sql.ErrConnDone)

	_, err = prodUcase.UserAdList(0, page)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_UserAdList_NoProduct(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

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
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

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

	//error
	prodRepo.EXPECT().SelectUserFavorite(uint64(0), gomock.Eq(page)).Return(nil, sql.ErrConnDone)

	_, err = prodUcase.GetUserFavorite(0, page)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_GetUserFavorite_NoProduct(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

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
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	prodRepo.EXPECT().InsertProductLike(uint64(0), uint64(0)).Return(nil)
	prodRepo.EXPECT().UpdateProductLikes(uint64(0), +1).Return(nil)

	err := prodUcase.LikeProduct(0, 0)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	prodRepo.EXPECT().InsertProductLike(uint64(0), uint64(0)).Return(sql.ErrConnDone)

	err = prodUcase.LikeProduct(0, 0)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_LikeProduct(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	prodRepo.EXPECT().InsertProductLike(uint64(0), uint64(0)).Return(pgx.PgError{Code: "23505"})

	err := prodUcase.LikeProduct(0, 0)
	assert.Equal(t, err, errors.UnexpectedInternal(pgx.PgError{Code: "23505"}))
}

func TestProductUsecase_DislikeProduct_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	prodRepo.EXPECT().DeleteProductLike(uint64(0), uint64(0)).Return(nil)
	prodRepo.EXPECT().UpdateProductLikes(uint64(0), -1).Return(nil)

	err := prodUcase.DislikeProduct(0, 0)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	prodRepo.EXPECT().DeleteProductLike(uint64(0), uint64(0)).Return(sql.ErrConnDone)

	err = prodUcase.DislikeProduct(0, 0)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_SetTariff_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	prodRepo.EXPECT().UpdateTariff(uint64(0), 0).Return(nil)

	err := prodUcase.SetTariff(0, 0)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	prodRepo.EXPECT().UpdateTariff(uint64(0), 0).Return(sql.ErrConnDone)

	err = prodUcase.SetTariff(0, 0)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_UpdatePhoto_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(prodTest, nil)
	uploadRepo.EXPECT().InsertPhotos(gomock.Eq([]*multipart.FileHeader{}), "static/product/").Return([]string{}, nil)
	prodRepo.EXPECT().InsertPhoto(gomock.Any()).Return(nil)
	//uploadRepo.EXPECT().RemovePhotos(gomock.Any()).Return(nil)

	prod, err := prodUcase.UpdatePhoto(prodTest.ID, uint64(0), []*multipart.FileHeader{})
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, prod, prodTest)

	//error
	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(prodTest, nil)
	uploadRepo.EXPECT().InsertPhotos(gomock.Eq([]*multipart.FileHeader{}), "static/product/").Return([]string{}, nil)
	prodRepo.EXPECT().InsertPhoto(gomock.Any()).Return(sql.ErrConnDone)
	//uploadRepo.EXPECT().RemovePhotos(gomock.Any()).Return(sql.ErrConnDone)

	_, err = prodUcase.UpdatePhoto(prodTest.ID, uint64(0), []*multipart.FileHeader{})
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))

	//another err
	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(prodTest, nil)
	uploadRepo.EXPECT().InsertPhotos(gomock.Eq([]*multipart.FileHeader{}), "static/product/").Return([]string{}, nil)
	prodRepo.EXPECT().InsertPhoto(gomock.Any()).Return(sql.ErrConnDone)

	_, err = prodUcase.UpdatePhoto(prodTest.ID, uint64(0), []*multipart.FileHeader{})
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_UpdatePhoto_WrongOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(prodTest, nil)

	_, err := prodUcase.UpdatePhoto(prodTest.ID, uint64(1), []*multipart.FileHeader{})
	assert.Equal(t, err, errors.Cause(errors.WrongOwner))
}

func TestProductUsecase_UpdatePhoto_NoProduct(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

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
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(prodTest, nil)
	uploadRepo.EXPECT().InsertPhotos(gomock.Any(), "static/product/").Return([]string{}, sql.ErrConnDone)

	_, err := prodUcase.UpdatePhoto(prodTest.ID, uint64(0), []*multipart.FileHeader{})
	assert.Equal(t, err.ErrorCode, errors.InternalError)
}

func TestProductUsecase_Close(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(&models.ProductData{OwnerID: prodTest.OwnerID}, nil)
	prodRepo.EXPECT().Close(&models.ProductData{OwnerID: prodTest.OwnerID}).Return(nil)

	err := prodUcase.Close(prodTest.ID, prodTest.OwnerID)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(nil, sql.ErrConnDone)
	err = prodUcase.Close(prodTest.ID, prodTest.OwnerID)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))

	//error
	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(&models.ProductData{OwnerID: prodTest.OwnerID+1}, nil)
	err = prodUcase.Close(prodTest.ID, prodTest.OwnerID)
	assert.Equal(t, err, errors.Cause(errors.WrongOwner))

	//error
	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(&models.ProductData{OwnerID: prodTest.OwnerID}, nil)
	prodRepo.EXPECT().Close(&models.ProductData{OwnerID: prodTest.OwnerID}).Return(sql.ErrConnDone)

	err = prodUcase.Close(prodTest.ID, prodTest.OwnerID)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_Edit(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	prodRepo.EXPECT().Update(&models.ProductData{OwnerID: prodTest.OwnerID}).Return(nil)

	err := prodUcase.Edit(&models.ProductData{OwnerID: prodTest.OwnerID})
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	prodRepo.EXPECT().Update(&models.ProductData{OwnerID: prodTest.OwnerID}).Return(sql.ErrConnDone)

	err = prodUcase.Edit(&models.ProductData{OwnerID: prodTest.OwnerID})
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_GetProduct(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(prodTest, nil)
	prodRepo.EXPECT().UpdateProductViews(prodTest.ID, 1).Return(nil)

	_, err := prodUcase.GetProduct(prodTest.ID)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(nil, sql.ErrConnDone)

	_, err = prodUcase.GetProduct(prodTest.ID)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))

	//error
	prodRepo.EXPECT().SelectByID(prodTest.ID).Return(prodTest, nil)
	prodRepo.EXPECT().UpdateProductViews(prodTest.ID, 1).Return(sql.ErrConnDone)

	_, err = prodUcase.GetProduct(prodTest.ID)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_TrendList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	trendsRepo.EXPECT().GetTrendsProducts(prodTest.OwnerID).Return([]uint64{}, nil)
	prodRepo.EXPECT().SelectTrands([]uint64{}, &prodTest.OwnerID).Return([]*models.ProductListData{}, nil)

	_, err := prodUcase.TrendList(&prodTest.OwnerID)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	trendsRepo.EXPECT().GetTrendsProducts(prodTest.OwnerID).Return([]uint64{}, nil)
	prodRepo.EXPECT().SelectTrands([]uint64{}, &prodTest.OwnerID).Return([]*models.ProductListData{}, sql.ErrConnDone)

	_, err = prodUcase.TrendList(&prodTest.OwnerID)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_GetProductReviewers(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)
	user := &models.UserData{
		ID:         0,
	}

	prodRepo.EXPECT().SelectProductReviewers(prodTest.ID, prodTest.OwnerID).Return([]*models.UserData{user}, nil)

	_, err := prodUcase.GetProductReviewers(prodTest.ID, prodTest.OwnerID)
	assert.Equal(t, err, (*errors.Error)(nil))

	//empty
	prodRepo.EXPECT().SelectProductReviewers(prodTest.ID, prodTest.OwnerID).Return([]*models.UserData{}, nil)

	_, err = prodUcase.GetProductReviewers(prodTest.ID, prodTest.OwnerID)
	assert.Equal(t, err, (*errors.Error)(nil))

	//err
	prodRepo.EXPECT().SelectProductReviewers(prodTest.ID, prodTest.OwnerID).Return(nil, sql.ErrConnDone)

	_, err = prodUcase.GetProductReviewers(prodTest.ID, prodTest.OwnerID)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_SetProductBuyer(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	prodRepo.EXPECT().InsertProductBuyer(prodTest.ID, prodTest.OwnerID).Return(nil)

	err := prodUcase.SetProductBuyer(prodTest.ID, prodTest.OwnerID)
	assert.Equal(t, err, (*errors.Error)(nil))

	//err
	prodRepo.EXPECT().InsertProductBuyer(prodTest.ID, prodTest.OwnerID).Return(sql.ErrConnDone)

	err = prodUcase.SetProductBuyer(prodTest.ID, prodTest.OwnerID)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_CreateProductReview(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)

	review := &models.Review{
		Content:        "dsd",
		Rating:         2,
		ReviewerID:     1,
		ProductID:      0,
		TargetID:       3,
		Type:           "seller",
	}

	prodRepo.EXPECT().CheckProductReview(review.ProductID, review.Type, review.ReviewerID).Return(false, nil)
	prodRepo.EXPECT().InsertReview(review).Return(nil)

	err := prodUcase.CreateProductReview(review)
	assert.Equal(t, err, (*errors.Error)(nil))

	//err
	prodRepo.EXPECT().CheckProductReview(review.ProductID, review.Type, review.ReviewerID).Return(false, sql.ErrConnDone)

	err = prodUcase.CreateProductReview(review)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))

	//err
	prodRepo.EXPECT().CheckProductReview(review.ProductID, review.Type, review.ReviewerID).Return(true, nil)

	err = prodUcase.CreateProductReview(review)
	assert.Equal(t, err, errors.Cause(errors.ReviewExist))

	//err
	prodRepo.EXPECT().CheckProductReview(review.ProductID, review.Type, review.ReviewerID).Return(false, nil)
	prodRepo.EXPECT().InsertReview(review).Return(sql.ErrConnDone)

	err = prodUcase.CreateProductReview(review)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_GetUserReviews(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)
	rv := &models.Review{ID: 0}

	prodRepo.EXPECT().SelectUserReviews(prodTest.OwnerID, "seller", &models.PageWithSort{}).Return([]*models.Review{rv}, nil)

	_, err := prodUcase.GetUserReviews(prodTest.OwnerID, "seller", &models.PageWithSort{})
	assert.Equal(t, err, (*errors.Error)(nil))

	//empty
	prodRepo.EXPECT().SelectUserReviews(prodTest.OwnerID, "seller", &models.PageWithSort{}).Return([]*models.Review{}, nil)

	_, err = prodUcase.GetUserReviews(prodTest.OwnerID, "seller", &models.PageWithSort{})
	assert.Equal(t, err, (*errors.Error)(nil))

	//err
	prodRepo.EXPECT().SelectUserReviews(prodTest.OwnerID, "seller", &models.PageWithSort{}).Return(nil, sql.ErrConnDone)

	_, err = prodUcase.GetUserReviews(prodTest.OwnerID, "seller", &models.PageWithSort{})
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestProductUsecase_GetWaitingReviews(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mock.NewMockProductRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	trendsRepo := tMock.NewMockTrendsRepository(ctrl)
	prodUcase := NewProductUsecase(prodRepo, uploadRepo, trendsRepo)
	rv := &models.WaitingReview{ProductID: 0}

	prodRepo.EXPECT().SelectWaitingReviews(prodTest.OwnerID, "seller", &models.Page{}).Return([]*models.WaitingReview{rv}, nil)

	_, err := prodUcase.GetWaitingReviews(prodTest.OwnerID, "seller", &models.Page{})
	assert.Equal(t, err, (*errors.Error)(nil))

	//empty
	prodRepo.EXPECT().SelectWaitingReviews(prodTest.OwnerID, "seller", &models.Page{}).Return([]*models.WaitingReview{}, nil)

	_, err = prodUcase.GetWaitingReviews(prodTest.OwnerID, "seller", &models.Page{})
	assert.Equal(t, err, (*errors.Error)(nil))

	//err
	prodRepo.EXPECT().SelectWaitingReviews(prodTest.OwnerID, "seller", &models.Page{}).Return(nil, sql.ErrConnDone)

	_, err = prodUcase.GetWaitingReviews(prodTest.OwnerID, "seller", &models.Page{})
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}
