package repository

import (
	"context"
	"database/sql"
	"strings"

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
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	query := tx.QueryRow(
		`
				INSERT INTO product(name, date, amount, description, category, owner_id)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id`,
		product.Name,
		product.Date,
		product.Amount,
		product.Description,
		product.Category,
		product.OwnerID)

	err = query.Scan(&product.ID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (pr *ProductRepository) SelectByID(productID uint64) (*models.ProductData, error) {
	product := &models.ProductData{}

	query := pr.dbConn.QueryRow(
		`
				SELECT p.id, p.name, p.date, p.amount, p.description, p.category, p.owner_id, u.name, u.surname, p.likes, p.views, array_agg(pi.img_link)
				FROM product AS p
				inner JOIN users as u ON p.owner_id=u.id and p.id=$1
				left join product_images as pi on pi.product_id=p.id
				GROUP BY p.id, u.name, u.surname`,
		productID)

	var linkStr string

	err := query.Scan(
		&product.ID,
		&product.Name,
		&product.Date,
		&product.Amount,
		&product.Description,
		&product.Category,
		&product.OwnerID,
		&product.OwnerName,
		&product.OwnerSurname,
		&product.Likes,
		&product.Views,
		&linkStr)
	if err != nil {
		return nil, err
	}
	linkStr = linkStr[1 : len(linkStr)-1]
	if linkStr != "NULL" {
		product.LinkImages = strings.Split(linkStr, ",")
	}
	return product, nil
}

func (pr *ProductRepository) SelectLatest(content *models.Content) ([]*models.ProductListData, error) {
	var products []*models.ProductListData

	query, err := pr.dbConn.Query(
		`
				SELECT p.id, p.name, p.date, p.amount, array_agg(pi.img_link)
				FROM product as p
				left join product_images as pi on pi.product_id=p.id
				GROUP BY p.id
				ORDER BY p.date DESC 
				LIMIT $1 OFFSET $2`,
		content.Count,
		content.From)
	if err != nil {
		return nil, err
	}

	defer query.Close()

	var linkStr string

	for query.Next() {
		product := &models.ProductListData{}

		err := query.Scan(
			&product.ID,
			&product.Name,
			&product.Date,
			&product.Amount,
			&linkStr)
		if err != nil {
			return nil, err
		}

		linkStr = linkStr[1 : len(linkStr)-1]
		if linkStr != "NULL" {
			product.LinkImages = strings.Split(linkStr, ",")
			products = append(products, product)
		}
	}

	if err := query.Err(); err != nil {
		return nil, err
	}
	return products, err

}

func (pr *ProductRepository) InsertPhoto(content *models.ProductData) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM product_images WHERE product_id=$1`,
		content.ID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	for _, photo := range content.LinkImages {
		_, err = tx.Exec(
			`INSERT INTO product_images(product_id, img_link)
		VALUES ($1, $2)`,
			content.ID, photo)

		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
