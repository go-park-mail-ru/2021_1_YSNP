package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category"
)

type CategoryRepository struct {
	dbConn *sql.DB
}

func NewCategoryRepository(conn *sql.DB) category.CategoryRepository {
	return &CategoryRepository{
		dbConn: conn,
	}
}

func (cat *CategoryRepository)  GetAllCategories() ([]*models.Category, error) {
	var categories []*models.Category

	query, err := cat.dbConn.Query(`SELECT title from category`)
	if err != nil {
		return nil, err
	}

	defer query.Close()

	for query.Next() {
		category := &models.Category{}

		err := query.Scan(&category.Title)

		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := query.Err(); err != nil {
		return nil, err
	}
	return categories, err
}