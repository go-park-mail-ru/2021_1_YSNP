package product

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

//go:generate mockgen -destination=./mocks/mock_product_repo.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product ProductRepository

type ProductRepository interface {
	Insert(product *models.ProductData) error
	SelectByID(productID uint64) (*models.ProductData, error)
	SelectLatest(content *models.Content) ([]*models.ProductListData, error)
	InsertPhoto(content *models.ProductData) error
}
