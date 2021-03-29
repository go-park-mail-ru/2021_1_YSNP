package usecase

import (
	"database/sql"
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
		return errors.UnexpectedInternal(err)
	}

	return nil
}

func (pu *ProductUsecase) GetByID(productID uint64) (*models.ProductData, *errors.Error) {
	product, err := pu.productRepo.SelectByID(productID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Cause(errors.ProductNotExist)
	case err != nil:
		return nil, errors.UnexpectedInternal(err)
	}

	return product, nil
}

func (pu *ProductUsecase) ListLatest(content *models.Content) ([]*models.ProductListData, *errors.Error) {
	products, err := pu.productRepo.SelectLatest(content)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	if len(products) == 0 {
		return []*models.ProductListData{}, nil
	}

	return products, nil
}

func (pu *ProductUsecase) UpdatePhoto(productID uint64, newPhoto []string) (*models.ProductData, *errors.Error) {
	product, errE := pu.GetByID(productID)
	if errE != nil {
		return nil, errE
	}

	oldPhotos := product.LinkImages
	product.LinkImages = newPhoto
	err := pu.productRepo.InsertPhoto(product)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	if len(oldPhotos) != 0 {
		for _, photo := range oldPhotos {
			err := os.Remove(photo)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}
		}
	}

	return product, nil
}

