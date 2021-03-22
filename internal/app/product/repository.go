package product

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

type ProductRepository interface {
	Insert(product *models.ProductData) error
	SelectByID(productID uint64) (*models.ProductData, error)
	SelectLatest() ([]*models.ProductListData, error)
}
