package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product"
)

type ProductRepository struct {
	dbConn *sql.DB
}

func NewProductRepository(conn *sql.DB) product.ProductRepository {
	return &ProductRepository{
		dbConn: conn,
	}
}

func (pr *ProductRepository) Insert(product *models.ProductData) error {
	panic("implement me")
}

func (pr *ProductRepository) SelectByID(productID uint64) (*models.ProductData, error) {
	panic("implement me")
}

func (pr *ProductRepository) SelectLatest() ([]*models.ProductListData, error) {
	panic("implement me")
}