package product

import (
	"mime/multipart"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

//go:generate mockgen -destination=./mocks/mock_product_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product ProductUsecase

type ProductUsecase interface {
	Create(product *models.ProductData) *errors.Error
	Close(productID uint64, ownerID uint64) *errors.Error
	Edit(product *models.ProductData) *errors.Error
	Delete(productID uint64) *errors.Error

	UpdatePhoto(productID uint64, ownerID uint64, filesHeaders []*multipart.FileHeader) (*models.ProductData, *errors.Error)
	SetTariff(productID uint64, tariff int) *errors.Error
	CreateProductReview(review *models.Review) *errors.Error

	GetProduct(productID uint64) (*models.ProductData, *errors.Error)
	ListLatest(userID *uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)
	UserAdList(userId uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)
	GetUserFavorite(userID uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)
	SetProductBuyer(productID uint64, buyerID uint64) *errors.Error
	GetProductReviewers(productID uint64, userID uint64) ([]*models.UserData, *errors.Error)
	GetUserReviews(userID uint64, reviewType string, content *models.PageWithSort) ([]*models.Review, *errors.Error)
	GetWaitingReviews(userID uint64, reviewType string, content *models.Page) ([]*models.WaitingReview, *errors.Error)

	LikeProduct(userID uint64, productID uint64) *errors.Error
	DislikeProduct(userID uint64, productID uint64) *errors.Error

	TrendList(userID *uint64) ([]*models.ProductListData, *errors.Error)
	RecommendationList(productID uint64, userID uint64) ([]*models.ProductListData, *errors.Error)

	GetByID(productID uint64) (*models.ProductData, *errors.Error)
}
