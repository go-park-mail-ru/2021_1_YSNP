package usecase

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product"
	"os"
	"time"
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
	product.Date = time.Now().UTC().Format("2006-01-02")

	err := pu.productRepo.Insert(product)
	if err != nil {
		//TODO: создать ошибку
		return &errors.Error{Message: err.Error()}
	}

	return nil
}

func (pu *ProductUsecase) GetByID(productID uint64) (*models.ProductData, *errors.Error) {
	product, err := pu.productRepo.SelectByID(productID)
	if err != nil {
		//TODO: создать ошибку
		return nil, &errors.Error{Message: err.Error()}
	}

	return product, nil
}

func (pu *ProductUsecase) ListLatest(content *models.Content) ([]*models.ProductListData, *errors.Error) {
	products, err := pu.productRepo.SelectLatest(content)
	if err != nil {
		//TODO: создать ошибку
		return nil, &errors.Error{Message: err.Error()}
	}

	return products, nil
}

func (pu *ProductUsecase) UpdatePhoto(productID uint64, newPhoto []string) (*models.ProductData, *errors.Error) {
	product, err := pu.productRepo.SelectByID(productID)
	if err != nil {
		//TODO: создать ошибку
		return nil, &errors.Error{Message: err.Error()}
	}

	oldPhotos := product.LinkImages
	product.LinkImages = newPhoto
	err = pu.productRepo.InsertPhoto(product)
	if err != nil {
		//TODO: создать ошибку
		return nil, &errors.Error{Message: err.Error()}
	}

	if len(oldPhotos) != 0 {
		for _, photo := range oldPhotos {
			err := os.Remove(photo)
			if err != nil {
				//TODO: создать ошибку
				return nil, &errors.Error{Message: err.Error()}
			}
		}
	}

	return product, nil
}

