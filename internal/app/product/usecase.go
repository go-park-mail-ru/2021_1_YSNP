package product

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

//go:generate mockgen -destination=./mocks/mock_product_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product ProductUsecase

type ProductUsecase interface {
	Create(product *models.ProductData) *errors.Error
	Close(product *models.ProductData, userID int) *errors.Error
	UpdatePhoto(productID uint64, newAvatar []string) (*models.ProductData, *errors.Error)
	SetTariff(productID uint64, tariff int) *errors.Error

	GetByID(productID uint64) (*models.ProductData, *errors.Error)
	ListLatest(userID *uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)
	UserAdList(userId uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)
	GetUserFavorite(userID uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)

	LikeProduct(userID uint64, productID uint64) *errors.Error
	DislikeProduct(userID uint64, productID uint64) *errors.Error
}
