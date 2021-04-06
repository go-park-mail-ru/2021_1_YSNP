package product

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

type ProductRepository interface {
	Insert(product *models.ProductData) error
	SelectByID(productID uint64) (*models.ProductData, error)
	SelectLatest(content *models.Page) ([]*models.ProductListData, error)
	SelectUserAd(userId uint64, content *models.Page) ([]*models.ProductListData, error)
	InsertPhoto(content *models.ProductData) error
}
