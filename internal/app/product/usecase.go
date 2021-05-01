package product

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	errors2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"mime/multipart"
)

//go:generate mockgen -destination=./mocks/mock_product_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product ProductUsecase

type ProductUsecase interface {
	Create(product *models.ProductData) *errors2.Error
	UpdatePhoto(productID uint64, ownerID uint64, filesHeaders []*multipart.FileHeader) (*models.ProductData, *errors2.Error)
	SetTariff(productID uint64, tariff int) *errors2.Error

	GetByID(productID uint64) (*models.ProductData, *errors2.Error)
	
	TrendList(userID *uint64) ([]*models.ProductListData, *errors2.Error)

	ListLatest(userID *uint64, content *models.Page) ([]*models.ProductListData, *errors2.Error)
	UserAdList(userId uint64, content *models.Page) ([]*models.ProductListData, *errors2.Error)
	GetUserFavorite(userID uint64, content *models.Page) ([]*models.ProductListData, *errors2.Error)

	LikeProduct(userID uint64, productID uint64) *errors2.Error
	DislikeProduct(userID uint64, productID uint64) *errors2.Error
}
