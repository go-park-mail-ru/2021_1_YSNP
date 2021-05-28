package repository

import (
	"context"
	"database/sql"
	"strconv"
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

func (pr *ProductRepository) Update(product *models.ProductData) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	query := tx.QueryRow(
		`
		UPDATE product set 
			name = $1,
			amount = $2,
			description = $3,
			category_id = (SELECT cat.id from category as cat where cat.title = $4),
			longitude = $5,
			latitude = $6,
			address = $7
		WHERE id = $8
		RETURNING id;
		`,
		product.Name,
		product.Amount,
		product.Description,
		product.Category,
		product.Longitude,
		product.Latitude,
		product.Address,
		product.ID,
	)

	err = query.Scan(&product.ID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}
	_, err = tx.Exec(
		`DELETE FROM product_images 
                WHERE product_id=$1`,
		product.ID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	for _, photo := range product.LinkImages {
		_, err = tx.Exec(
			`INSERT INTO product_images(product_id, img_link)
		            VALUES ($1, $2)`,
			product.ID,
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

func (pr *ProductRepository) Close(product *models.ProductData) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`Update product 
                set close = true 
                where id = $1`,
		product.ID)
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

func (pr *ProductRepository) InsertPhoto(content *models.ProductData) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	/*_, err = tx.Exec(
			`DELETE FROM product_images
	                WHERE product_id=$1`,
			content.ID)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return rollbackErr
			}

			return err
		}*/

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
	query := pr.dbConn.QueryRow(
		`
				SELECT p.id, p.name, p.date, p.amount, p.description, cat.title, p.owner_id, u.name, u.surname, u.avatar, u.score, u.reviews, p.likes, p.views, p.longitude, p.latitude, p.address, array_agg(pi.img_link), p.tariff, p.close
				FROM product AS p
				inner JOIN users as u ON p.owner_id=u.id and p.id=$1
				left join product_images as pi on pi.product_id=p.id
				left join category as cat on cat.id=p.category_id
				GROUP BY p.id, cat.title, u.name, u.surname, u.avatar, u.score, u.reviews`,
		productID)

	product := &models.ProductData{}
	var linkStr string
	var date time.Time
	var reviews int

	err := query.Scan(
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
		&product.OwnerRating,
		&reviews,
		&product.Likes,
		&product.Views,
		&product.Longitude,
		&product.Latitude,
		&product.Address,
		&linkStr,
		&product.Tariff,
		&product.Close)
	if err != nil {
		return nil, err
	}

	product.Date = date.Format("2006-01-02")
	if reviews != 0 {
		product.OwnerRating = product.OwnerRating / float64(reviews)
	}

	linkStr = linkStr[1 : len(linkStr)-1]
	if linkStr != "NULL" {
		product.LinkImages = strings.Split(linkStr, ",")
	}

	return product, nil
}

func (pr *ProductRepository) SelectTrands(idArray []uint64, userID *uint64) ([]*models.ProductListData, error) {
	var products []*models.ProductListData

	queries := `SELECT p.id, p.name, p.date, p.amount, array_agg(pi.img_link), uf.user_id, p.tariff
	FROM product as p
	left join product_images as pi on pi.product_id=p.id
	left join user_favorite uf on p.id = uf.product_id and uf.user_id = $1
	WHERE p.close=false and p.id IN (`

	var val []interface{}
	val = append(val, *userID)

	hasId := false
	for i, item := range idArray {
		queries += " $" + strconv.Itoa(i+2) + ","
		val = append(val, item)
		hasId = true
	}

	if hasId {
		queries = queries[:len(queries)-1]
	} else {
		queries += "null"
	}

	queries += `)
			GROUP BY p.id, uf.user_id
			ORDER BY p.date DESC 
		`
	query, err := pr.dbConn.Query(
		queries,
		val...,
	)
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

func (pr *ProductRepository) SelectLatest(userID *uint64, content *models.Page) ([]*models.ProductListData, error) {
	var products []*models.ProductListData

	query, err := pr.dbConn.Query(
		`
				SELECT p.id, p.name, p.date, p.amount, array_agg(pi.img_link), uf.user_id, p.tariff
				FROM product as p
				left join product_images as pi on pi.product_id=p.id
				left join user_favorite uf on p.id = uf.product_id and uf.user_id = $3
				WHERE p.close=false
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
				SELECT p.id, p.name, p.date, p.amount, array_agg(pi.img_link), p.tariff, p.close
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
			&product.Tariff,
			&product.Close)

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

func (pr *ProductRepository) UpdateProductLikes(productID uint64, count int) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`
				UPDATE product
                SET likes=likes + $1
                WHERE product.id = $2`,
		count,
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

func (pr *ProductRepository) UpdateProductViews(productID uint64, count int) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`
				UPDATE product
                SET views=views + $1
                WHERE product.id = $2`,
		count,
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

func (pr *ProductRepository) SelectProductReviewers(productID uint64, userID uint64) ([]*models.UserData, error) {
	var users []*models.UserData

	query, err := pr.dbConn.Query(
		`
				SELECT u.id, u.name, u.avatar
From users as u
left join user_chats as uc on uc.product_id = $1 and uc.user_id = u.id and uc.user_id != $2
left join user_favorite as uf on uf.product_id = $1 and uf.user_id = u.id
where uf.user_id notnull or uc.chat_id notnull`,
		productID,
		userID)
	if err != nil {
		return nil, err
	}

	defer query.Close()

	for query.Next() {
		user := &models.UserData{}

		err := query.Scan(
			&user.ID,
			&user.Name,
			&user.LinkImages)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := query.Err(); err != nil {
		return nil, err
	}

	return users, err
}

func (pr *ProductRepository) InsertProductBuyer(productID uint64, buyerID uint64) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`
				UPDATE product
                SET buyer_id=$1, buyer_left_review = false, seller_left_review = false
                WHERE product.id = $2`,
		buyerID,
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

func (pr *ProductRepository) InsertReview(review *models.Review) error {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	query := tx.QueryRow(
		`
				INSERT INTO reviews (creation_time) 
			VALUES ($1) 
			RETURNING id, creation_time; `,
		time.Now())

	err = query.Scan(&review.ID, &review.CreationTime)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	_, err = tx.Exec(
		`
				INSERT INTO user_reviews (review_id, content, rating, reviewer_id, product_id,  target_id, type)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		review.ID,
		review.Content,
		review.Rating,
		review.ReviewerID,
		review.ProductID,
		review.TargetID,
		review.Type)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	updateQuery := `UPDATE product `
	if review.Type == "buyer" {
		updateQuery += `SET buyer_left_review = true
                WHERE product.id = $1`
	} else {
		updateQuery += `SET seller_left_review = true
                WHERE product.id = $1`
	}
	_, err = tx.Exec(updateQuery,
		review.ProductID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	_, err = tx.Exec(`UPDATE users SET score= score + $1, reviews = reviews + 1, new_revs = new_revs + 1 WHERE id = $2`,
		int64(review.Rating), review.TargetID)
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

func (pr *ProductRepository) CheckProductReview(productID uint64, reviewType string, reviewerID uint64) (bool, error) {
	var result bool
	var revID uint64
	selectQuery := `SELECT `
	if reviewType == "buyer" {
		selectQuery += `buyer_id, buyer_left_review FROM product
                WHERE product.id = $1`
	} else {
		selectQuery += `owner_id, seller_left_review FROM product
                WHERE product.id = $1`
	}
	err := pr.dbConn.QueryRow(selectQuery,
		productID).Scan(&revID, &result)
	if err != nil {
		return false, err
	}
	if revID != reviewerID {
		return true, nil
	}

	return result, nil
}

func (pr *ProductRepository) SelectUserReviews(userID uint64, reviewType string, content *models.PageWithSort, loggedUser int64) ([]*models.Review, error) {
	tx, err := pr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	var reviews []*models.Review

	selectQuery := `WITH ORDERED AS
(SELECT r.id, r.creation_time, ur.content, ur.rating, ur.type, ur.reviewer_id, u.name as uname, u.avatar, ur.product_id, p.name, pi.img_link, ROW_NUMBER() OVER (PARTITION BY r.id) As rn
				FROM user_reviews AS ur
				JOIN reviews AS r on ur.review_id = r.id
				JOIN users AS u ON u.id = ur.reviewer_id
				JOIN product AS p ON p.id = ur.product_id
				Left Join product_images pi on p.id = pi.product_id
WHERE ur.target_id = $1 and ur.type = $2`
	if content.Sort == "date" {
		selectQuery += `
			ORDER BY r.creation_time desc
		`
	}
	if content.Sort == "rate" {
		selectQuery += `
			ORDER BY ur.rating desc
		`
	}

	selectQuery += `
			)
		SELECT
		id, creation_time, content, rating, type, reviewer_id, uname, avatar, product_id, name, img_link
		FROM
			ORDERED
		WHERE
			rn = 1
			LIMIT $3 OFFSET $4`

	var queryData []interface{}
	queryData = append(queryData, userID)
	queryData = append(queryData, reviewType)
	queryData = append(queryData, content.Count)
	queryData = append(queryData, content.From*content.Count)

	query, err := tx.Query(selectQuery, queryData...)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	defer query.Close()

	for query.Next() {
		review := &models.Review{}

		err := query.Scan(
			&review.ID,
			&review.CreationTime,
			&review.Content,
			&review.Rating,
			&review.Type,
			&review.ReviewerID,
			&review.ReviewerName,
			&review.ReviewerAvatar,
			&review.ProductID,
			&review.ProductName,
			&review.ProductImage)

		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return nil, rollbackErr
			}
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := query.Err(); err != nil {
		return nil, err
	}

	if (loggedUser == int64(userID)) {
		_, err = tx.Exec(`UPDATE users SET new_revs = 0 WHERE id = $1`, userID)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return nil, rollbackErr
			}
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return reviews, err
}

func (pr *ProductRepository) SelectWaitingReviews(userID uint64, reviewType string, content *models.Page) ([]*models.WaitingReview, error) {
	var reviews []*models.WaitingReview

	selectQuery := `
WITH ORDERED AS
(SELECT u.id as uid, u.name as uname, u.avatar, p.id as pid, p.name as pname, p.owner_id, pi.img_link, ROW_NUMBER() OVER (PARTITION BY p.id) As rn
				FROM product as p
				left JOIN users AS u ON `

	if reviewType == "seller" {
		selectQuery += `
u.id = p.owner_id and p.buyer_left_review = false
				Left Join product_images pi on p.id = pi.product_id
WHERE p.buyer_id = $1 and u.id notnull
    ORDER BY p.date DESC
)
SELECT
uid, uname, avatar, pid, pname, owner_id, img_link
FROM
    ORDERED
WHERE
    rn = 1
    LIMIT $2 OFFSET $3`
	} else {
		selectQuery += `
u.id = p.buyer_id and p.seller_left_review = false
				Left Join product_images pi on p.id = pi.product_id
WHERE p.owner_id = $1 and u.id notnull
    ORDER BY p.date DESC
)
SELECT
uid, uname, avatar, pid, pname, owner_id, img_link
FROM
    ORDERED
WHERE
    rn = 1
    LIMIT $2 OFFSET $3`
	}
	query, err := pr.dbConn.Query(
		selectQuery,
		userID,
		content.Count,
		content.From*content.Count)
	if err != nil {
		return nil, err
	}

	defer query.Close()

	for query.Next() {
		var owner uint64
		review := &models.WaitingReview{}

		err := query.Scan(
			&review.TargetID,
			&review.TargetName,
			&review.TargetAvatar,
			&review.ProductID,
			&review.ProductName,
			&owner,
			&review.ProductImage)

		if err != nil {
			return nil, err
		}

		if owner == userID {
			review.Type = "seller"
		} else {
			review.Type = "buyer"
		}
		reviews = append(reviews, review)
	}

	if err := query.Err(); err != nil {
		return nil, err
	}

	return reviews, err
}
