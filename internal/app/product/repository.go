package product

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

//go:generate mockgen -destination=./mocks/mock_product_repo.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product ProductRepository

type ProductRepository interface {
	Insert(product *models.ProductData) error
	InsertPhoto(content *models.ProductData) error
	UpdateTariff(productID uint64, tariff int) error

	SelectByID(productID uint64) (*models.ProductData, error)
	SelectLatest(userID *uint64, content *models.Page) ([]*models.ProductListData, error)
	SelectUserAd(userId uint64, content *models.Page) ([]*models.ProductListData, error)
	SelectUserFavorite(userID uint64, content *models.Page) ([]*models.ProductListData, error)

	InsertProductLike(userID uint64, productID uint64) error
	DeleteProductLike(userID uint64, productID uint64) error
}
