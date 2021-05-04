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
  
	UpdatePhoto(productID uint64, ownerID uint64, filesHeaders []*multipart.FileHeader) (*models.ProductData, *errors.Error)
	SetTariff(productID uint64, tariff int) *errors.Error

	GetProduct(productID uint64) (*models.ProductData, *errors.Error)
  TrendList(userID *uint64) ([]*models.ProductListData, *errors.Error)
	ListLatest(userID *uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)
	UserAdList(userId uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)
	GetUserFavorite(userID uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)

	LikeProduct(userID uint64, productID uint64) *errors.Error
	DislikeProduct(userID uint64, productID uint64) *errors.Error

	GetByID(productID uint64) (*models.ProductData, *errors.Error)
}
