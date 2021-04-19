package product

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"mime/multipart"
)

//go:generate mockgen -destination=./mocks/mock_product_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product ProductUsecase

type ProductUsecase interface {
	Create(product *models.ProductData) *errors.Error
	UpdatePhoto(productID uint64, ownerID uint64, files []*multipart.FileHeader) *errors.Error
	SetTariff(productID uint64, tariff int) *errors.Error

	GetByID(productID uint64) (*models.ProductData, *errors.Error)
	ListLatest(userID *uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)
	UserAdList(userId uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)
	GetUserFavorite(userID uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)

	LikeProduct(userID uint64, productID uint64) *errors.Error
	DislikeProduct(userID uint64, productID uint64) *errors.Error
}
