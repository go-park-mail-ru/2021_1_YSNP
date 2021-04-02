package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search"
)

type SearchRepository struct {
	dbConn *sql.DB
}

func (s SearchRepository) SelectByFilter(data *models.Search) ([]*models.ProductData, error) {
	var products []*models.ProductData

	query, err := s.dbConn.Query("SELECT id, name, date, amount from product where name like $1","%" + data.Search + "%" )
	if err != nil {
		return nil, err
	}

	defer query.Close()

	//var linkStr string

	for query.Next() {
		product := &models.ProductData{}

		err := query.Scan(
			&product.ID,
			&product.Name,
			&product.Date,
			&product.Amount,
			)
		if err != nil {
			return nil, err
		}


		products = append(products, product)
/*
		linkStr = linkStr[1 : len(linkStr)-1]
		if linkStr != "NULL" {
			product.LinkImages = strings.Split(linkStr, ",")
		}*/
	}

	if err := query.Err(); err != nil {
		return nil, err
	}
	return products, err
}

func NewProductRepository(conn *sql.DB) search.SearchRepository {
	return & SearchRepository{
		dbConn: conn,
	}
}