package usecase

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product"
)

type ProductUsecase struct {
	productRepo product.ProductRepository
}

func NewProductUsecase(repo product.ProductRepository) product.ProductUsecase {
	return &ProductUsecase{
		productRepo: repo,
	}
}

func (pu *ProductUsecase) Create(product *models.ProductData) *errors.Error {
	err := pu.productRepo.Insert(product)
	if err != nil {
		//TODO: создать ошибку
	}

	return nil
}

func (pu *ProductUsecase) GetByID(productID uint64) (*models.ProductData, *errors.Error) {
	product, err := pu.productRepo.SelectByID(productID)
	if err != nil {
		//TODO: создать ошибку
	}

	return product, nil
}

func (pu *ProductUsecase) ListLatest() ([]*models.ProductListData, *errors.Error) {
	products, err := pu.productRepo.SelectLatest()
	if err != nil {
		//TODO: создать ошибку
	}

	return products, nil
}
