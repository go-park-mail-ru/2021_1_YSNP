package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category"
)

type CategoryRepository struct {
	dbConn *sql.DB
}

func NewProductRepository(conn *sql.DB) category.CategoryRepository {
	return &CategoryRepository{
		dbConn: conn,
	}
}

func (cat *CategoryRepository)  GetCategory() ([]*models.Category, error) {
	return nil, nil
}