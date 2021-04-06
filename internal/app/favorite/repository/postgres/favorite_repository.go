package postgres

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/favorite"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"strings"
	"time"
)

type FavoriteRepository struct {
	dbConn *sql.DB
}

func NewFavoriteRepository(conn *sql.DB) favorite.FavoriteRepository {
	return &FavoriteRepository{
		dbConn: conn,
	}
}

func (pr *FavoriteRepository) InsertProduct(userID uint64, productID uint64) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`
				INSERT INTO user_favorite
                (user_id, product_id)
                VALUES ($1, $2) `,
		userID,
		productID)
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

func (pr *FavoriteRepository) DeleteProduct(userID uint64, productID uint64) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`
				DELETE from user_favorite
                where user_id=$1 and product_id=$2`,
		userID,
		productID)
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

func (pr *FavoriteRepository) SelectUserFavorite(userID uint64, content *models.Page) ([]*models.ProductListData, error) {
	var products []*models.ProductListData

	query, err := pr.dbConn.Query(
		`
				SELECT p.id, p.name, p.date, p.amount, array_agg(pi.img_link)
                FROM user_favorite
                JOIN product p ON p.id = user_favorite.product_id
                LEFT JOIN product_images AS pi ON pi.product_id = p.id
                WHERE user_id=$1
                GROUP BY p.id
                ORDER BY p.date DESC
                LIMIT $2 OFFSET $3`,
		userID,
		content.Count,
		content.From*content.Count)
	if err != nil {
		return nil, err
	}

	defer query.Close()

	var linkStr string
	var date time.Time

	for query.Next() {
		product := &models.ProductListData{}

		err := query.Scan(
			&product.ID,
			&product.Name,
			&date,
			&product.Amount,
			&linkStr)

		if err != nil {
			return nil, err
		}

		product.Date = date.Format("2006-01-02")
		linkStr = linkStr[1 : len(linkStr)-1]
		if linkStr != "NULL" {
			product.LinkImages = strings.Split(linkStr, ",")
		}
		products = append(products, product)
	}

	if err := query.Err(); err != nil {
		return nil, err
	}
	return products, err
}
