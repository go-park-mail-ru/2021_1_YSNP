package product

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

type ProductUsecase interface {
	Create(product *models.ProductData) *errors.Error
	UpdatePhoto(productID uint64, newAvatar []string) (*models.ProductData, *errors.Error)

	GetByID(productID uint64) (*models.ProductData, *errors.Error)
	ListLatest(content *models.Page) ([]*models.ProductListData, *errors.Error)
	ListAuthLatest(userID uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)
	UserAdList(userId uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)
	GetUserFavorite(userID uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)

	LikeProduct(userID uint64, productID uint64) *errors.Error
	DislikeProduct(userID uint64, productID uint64) *errors.Error
}
