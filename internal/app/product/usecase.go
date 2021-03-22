package product

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

type ProductUsecase interface {
	Create(product *models.ProductData) *errors.Error
	GetByID(productID uint64) (*models.ProductData, *errors.Error)
	ListLatest() ([]*models.ProductListData, *errors.Error)
}