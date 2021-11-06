package product

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

//go:generate mockgen -destination=./mocks/mock_product_repo.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product ProductRepository

type ProductRepository interface {
	Insert(product *models.ProductData) error
	Close(product *models.ProductData) error
	Update(product *models.ProductData) error
	Delete(productID uint64) error

	InsertPhoto(content *models.ProductData) error
	UpdateTariff(productID uint64, tariff int) error

	InsertProductBuyer(productID uint64, buyerID uint64) error
	InsertReview(review *models.Review) error
	CheckProductReview(productID uint64, reviewType string, reviewerID uint64) (bool, error)
	SelectUserReviews(userID uint64, reviewType string, content *models.PageWithSort) ([]*models.Review, error)

	SelectByID(productID uint64) (*models.ProductData, error)
	SelectTrands(idArray []uint64, userID *uint64) ([]*models.ProductListData, error)
	SelectLatest(userID *uint64, content *models.Page) ([]*models.ProductListData, error)
	SelectUserAd(userId uint64, content *models.Page) ([]*models.ProductListData, error)
	SelectUserFavorite(userID uint64, content *models.Page) ([]*models.ProductListData, error)
	SelectProductReviewers(productID uint64, userID uint64) ([]*models.UserData, error)
	SelectWaitingReviews(userID uint64, reviewType string, content *models.Page) ([]*models.WaitingReview, error)

	InsertProductLike(userID uint64, productID uint64) error
	DeleteProductLike(userID uint64, productID uint64) error

	UpdateProductLikes(productID uint64, count int) error
	UpdateProductViews(productID uint64, count int) error
}
