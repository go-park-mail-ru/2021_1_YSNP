package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

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
				INSERT INTO product(name, date, amount, description, category_id, owner_id, longitude, latitude, address)
				VALUES ($1, $2, $3, $4, (SELECT cat.id from category as cat where cat.title = $5), $6, $7, $8, $9)
				RETURNING id`,
		product.Name,
		product.Date,
		product.Amount,
		product.Description,
		product.Category,
		product.OwnerID,
		product.Longitude,
		product.Latitude,
		product.Address)

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

func(pr *ProductRepository)  Close(product *models.ProductData, userID int) error {
	err := pr.dbConn.QueryRow(`
	SELECT p.owner_id
	FROM product AS p
	WHERE p.id = $1`, product.ID).Scan(&product.OwnerID)

	if err != nil {
		return err
	}
	if int(product.OwnerID) !=  userID {
		return errors.New("userError")
	}

	tx, errtx := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if errtx != nil {
		return err
	}
	_, errU := tx.Exec(
		`Update product set close = true where id = $1`,
		product.ID)
	if errU != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	errtx = tx.Commit()
	if errtx != nil {
		return errtx
	}

	return nil
}

func (pr *ProductRepository) InsertPhoto(content *models.ProductData) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM product_images 
                WHERE product_id=$1`,
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
			content.ID,
			photo)
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

func (pr *ProductRepository) SelectByID(productID uint64) (*models.ProductData, error) {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	query := tx.QueryRow(
		`
				SELECT p.id, p.name, p.date, p.amount, p.description, cat.title, p.owner_id, u.name, u.surname, u.avatar, p.likes, p.views, p.longitude, p.latitude, p.address, array_agg(pi.img_link), p.tariff, p.close
				FROM product AS p
				inner JOIN users as u ON p.owner_id=u.id and p.id=$1
				left join product_images as pi on pi.product_id=p.id
				left join category as cat on cat.id=p.category_id
				GROUP BY p.id, cat.title, u.name, u.surname, u.avatar`,
		productID)

	product := &models.ProductData{}
	var linkStr string
	var date time.Time

	err = query.Scan(
		&product.ID,
		&product.Name,
		&date,
		&product.Amount,
		&product.Description,
		&product.Category,
		&product.OwnerID,
		&product.OwnerName,
		&product.OwnerSurname,
		&product.OwnerLinkImages,
		&product.Likes,
		&product.Views,
		&product.Longitude,
		&product.Latitude,
		&product.Address,
		&linkStr,
		&product.Tariff,
		&product.Close)

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}

		return nil, err
	}

	_, err = tx.Exec(
		`
				UPDATE product
                SET views=views + 1
                WHERE product.id = $1`,
		productID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}

		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	product.Date = date.Format("2006-01-02")
	linkStr = linkStr[1 : len(linkStr)-1]
	if linkStr != "NULL" {
		product.LinkImages = strings.Split(linkStr, ",")
	}

	return product, nil
}

func (pr *ProductRepository) SelectLatest(userID *uint64, content *models.Page) ([]*models.ProductListData, error) {
	var products []*models.ProductListData

	query, err := pr.dbConn.Query(
		`
				SELECT p.id, p.name, p.date, p.amount, array_agg(pi.img_link), uf.user_id, p.tariff
				FROM product as p
				left join product_images as pi on pi.product_id=p.id
				left join user_favorite uf on p.id = uf.product_id and uf.user_id = $3
				GROUP BY p.id, uf.user_id
				ORDER BY p.date DESC 
				LIMIT $1 OFFSET $2`,
		content.Count,
		content.From*content.Count,
		*userID)
	if err != nil {
		return nil, err
	}

	defer query.Close()

	var linkStr string
	var date time.Time

	for query.Next() {
		product := &models.ProductListData{}
		var user sql.NullInt64

		err := query.Scan(
			&product.ID,
			&product.Name,
			&date,
			&product.Amount,
			&linkStr,
			&user,
			&product.Tariff)
		if err != nil {
			return nil, err
		}

		product.UserLiked = false
		if userID != nil && user.Valid && uint64(user.Int64) == *userID {
			product.UserLiked = true
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

func (pr *ProductRepository) SelectUserAd(userId uint64, content *models.Page) ([]*models.ProductListData, error) {
	var products []*models.ProductListData

	query, err := pr.dbConn.Query(
		`
				SELECT p.id, p.name, p.date, p.amount, array_agg(pi.img_link), p.tariff
				FROM product as p
				left join product_images as pi on pi.product_id=p.id
				WHERE owner_id=$1
				GROUP BY p.id
				ORDER BY p.date DESC
				LIMIT $2 OFFSET $3`,
		userId,
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
			&linkStr,
			&product.Tariff)

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

func (pr *ProductRepository) SelectUserFavorite(userID uint64, content *models.Page) ([]*models.ProductListData, error) {
	var products []*models.ProductListData

	query, err := pr.dbConn.Query(
		`
				SELECT p.id, p.name, p.date, p.amount, array_agg(pi.img_link), p.tariff
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
			&linkStr,
			&product.Tariff)

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

func (pr *ProductRepository) InsertProductLike(userID uint64, productID uint64) error {
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

	_, err = tx.Exec(
		`
				UPDATE product
                SET likes=likes + 1
                WHERE product.id = $1`,
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

func (pr *ProductRepository) DeleteProductLike(userID uint64, productID uint64) error {
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

	_, err = tx.Exec(
		`
				UPDATE product
                SET likes=likes - 1
                WHERE product.id = $1`,
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

func (pr *ProductRepository) UpdateTariff(productID uint64, tariff int) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE product SET tariff=$1 WHERE id=$2`,
		tariff,
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
