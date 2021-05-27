package repository

import (
	"database/sql"
	"database/sql/driver"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

var prodTest = &models.ProductData{
	ID:          0,
	Name:        "tovar",
	Date:        "2010-02-02",
	Amount:      10000,
	Description: "Description product aaaaa",
	Category:    "0",
	OwnerID:     0,
	Longitude:   1,
	Latitude:    1,
	Address:     "Address",
	LinkImages:  []string{"test_str.jpg"},
	Tariff:      0,
	Close: false,
	OwnerRating: 0,
}

func TestProductRepository_Insert_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	mock.ExpectBegin()
	answer := sqlmock.NewRows([]string{"id"}).AddRow(prodTest.ID)
	mock.ExpectQuery(`INSERT INTO product`).WithArgs(
		prodTest.Name,
		prodTest.Date,
		prodTest.Amount,
		prodTest.Description,
		prodTest.Category,
		prodTest.OwnerID,
		prodTest.Longitude,
		prodTest.Latitude,
		prodTest.Address).WillReturnRows(answer)
	mock.ExpectCommit()

	err = prodRepo.Insert(prodTest)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductRepository_Insert_Error(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	mock.ExpectBegin()
	answer := sqlmock.NewRows([]string{"id"}).AddRow(prodTest.Name)
	mock.ExpectQuery(`INSERT INTO product`).WithArgs(
		prodTest.Name,
		prodTest.Date,
		prodTest.Amount,
		prodTest.Description,
		prodTest.Category,
		prodTest.OwnerID,
		prodTest.Longitude,
		prodTest.Latitude,
		prodTest.Address).WillReturnRows(answer)
	mock.ExpectRollback()

	err = prodRepo.Insert(prodTest)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductRepository_SelectByID_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	layout := "2006-01-02"
	date, _ := time.Parse(layout, prodTest.Date)

	linkStr := "{" + strings.Join(prodTest.LinkImages, ",") + "}"

	rows := sqlmock.NewRows([]string{"p.id", "p.name", "p.date", "p.amount", "p.description", "cat.title", "p.owner_id", "u.name",
		"u.surname", "u.avatar", "u.score", "u.reviews", "p.likes", "p.views", "p.longitude", "p.latitude", "p.address", "array_agg(pi.img_link)", "p.tariff", "p.close"})
	rows.AddRow(
		&prodTest.ID,
		&prodTest.Name,
		&date,
		&prodTest.Amount,
		&prodTest.Description,
		&prodTest.Category,
		&prodTest.OwnerID,
		&prodTest.OwnerName,
		&prodTest.OwnerSurname,
		&prodTest.OwnerLinkImages,
		&prodTest.OwnerRating,
		0,
		&prodTest.Likes,
		&prodTest.Views,
		&prodTest.Longitude,
		&prodTest.Latitude,
		&prodTest.Address,
		&linkStr,
		&prodTest.Tariff,
		&prodTest.Close)
	mock.ExpectQuery(`SELECT`).WithArgs(prodTest.ID).WillReturnRows(rows)
	//mock.ExpectExec(`UPDATE product`).WithArgs(prodTest.ID).WillReturnResult(sqlmock.NewResult(int64(prodTest.ID), 1))
	//mock.ExpectCommit()

	user, err := prodRepo.SelectByID(prodTest.ID)
	assert.Equal(t, prodTest, user)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductRepository_SelectByID_SelectErr(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	linkStr := "{" + strings.Join(prodTest.LinkImages, ",") + "}"

	rows := sqlmock.NewRows([]string{"p.id", "p.name", "p.date", "p.amount", "p.description", "cat.title", "p.owner_id", "u.name",
		"u.surname", "u.avatar", "p.likes", "p.views", "p.longitude", "p.latitude", "p.address", "array_agg(pi.img_link)", "p.tariff"})
	rows.AddRow(
		&prodTest.ID,
		&prodTest.Name,
		&prodTest.Date, //wrong date type
		&prodTest.Amount,
		&prodTest.Description,
		&prodTest.Category,
		&prodTest.OwnerID,
		&prodTest.OwnerName,
		&prodTest.OwnerSurname,
		&prodTest.OwnerLinkImages,
		&prodTest.Likes,
		&prodTest.Views,
		&prodTest.Longitude,
		&prodTest.Latitude,
		&prodTest.Address,
		&linkStr,
		&prodTest.Tariff)
	mock.ExpectQuery(`SELECT`).WithArgs(prodTest.ID).WillReturnRows(rows)

	_, err = prodRepo.SelectByID(prodTest.ID)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}


func TestProductRepository_SelectLatest_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	page := &models.Page{
		From:  1,
		Count: 10,
	}

	var userID uint64 = 1

	layout := "2006-01-02"
	date, _ := time.Parse(layout, prodTest.Date)

	linkStr := "{" + strings.Join(prodTest.LinkImages, ",") + "}"

	list := &models.ProductListData{
		ID:         0,
		Name:       "tovar",
		Date:       "2010-02-02",
		Amount:     10000,
		LinkImages: []string{"test_str.jpg"},
		UserLiked:  true,
		Tariff:     0,
	}

	rows := sqlmock.NewRows([]string{"p.id", "p.name", "p.date", "p.amount", "array_agg(pi.img_link)", "uf.user_id", "p.tariff"})
	rows.AddRow(
		&prodTest.ID,
		&prodTest.Name,
		&date,
		&prodTest.Amount,
		&linkStr,
		&userID,
		&prodTest.Tariff)
	mock.ExpectQuery(`SELECT`).WithArgs(10, 10, &userID).WillReturnRows(rows)

	user, err := prodRepo.SelectLatest(&userID, page)
	assert.Equal(t, list, user[0])
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductRepository_SelectUserAd_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	page := &models.Page{
		From:  1,
		Count: 10,
	}

	var userID uint64 = 1

	layout := "2006-01-02"
	date, _ := time.Parse(layout, prodTest.Date)

	linkStr := "{" + strings.Join(prodTest.LinkImages, ",") + "}"

	list := &models.ProductListData{
		ID:         0,
		Name:       "tovar",
		Date:       "2010-02-02",
		Amount:     10000,
		LinkImages: []string{"test_str.jpg"},
		UserLiked:  false,
		Tariff:     0,
		Close: false,
	}

	rows := sqlmock.NewRows([]string{"p.id", "p.name", "p.date", "p.amount", "array_agg(pi.img_link)", "p.tariff", "p.close"})
	rows.AddRow(
		&prodTest.ID,
		&prodTest.Name,
		&date,
		&prodTest.Amount,
		&linkStr,
		&prodTest.Tariff,
		&prodTest.Close)
	mock.ExpectQuery(`SELECT`).WithArgs(userID, 10, 10).WillReturnRows(rows)

	user, err := prodRepo.SelectUserAd(userID, page)
	assert.Equal(t, list, user[0])
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductRepository_SelectUserFavorite_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	page := &models.Page{
		From:  1,
		Count: 10,
	}

	var userID uint64 = 1

	layout := "2006-01-02"
	date, _ := time.Parse(layout, prodTest.Date)

	linkStr := "{" + strings.Join(prodTest.LinkImages, ",") + "}"

	list := &models.ProductListData{
		ID:         0,
		Name:       "tovar",
		Date:       "2010-02-02",
		Amount:     10000,
		LinkImages: []string{"test_str.jpg"},
		UserLiked:  false,
		Tariff:     0,
	}

	rows := sqlmock.NewRows([]string{"p.id", "p.name", "p.date", "p.amount", "array_agg(pi.img_link)", "p.tariff"})
	rows.AddRow(
		&prodTest.ID,
		&prodTest.Name,
		&date,
		&prodTest.Amount,
		&linkStr,
		&prodTest.Tariff)
	mock.ExpectQuery(`SELECT`).WithArgs(userID, 10, 10).WillReturnRows(rows)

	user, err := prodRepo.SelectUserFavorite(userID, page)
	assert.Equal(t, list, user[0])
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductRepository_InsertProductLike_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	var userID uint64 = 1

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO user_favorite`).WithArgs(
		userID,
		prodTest.ID).WillReturnResult(driver.ResultNoRows)
	//mock.ExpectExec(`UPDATE product`).WithArgs(prodTest.ID).WillReturnResult(sqlmock.NewResult(int64(prodTest.ID), 1))
	mock.ExpectCommit()

	err = prodRepo.InsertProductLike(userID, prodTest.ID)
	assert.NoError(t, err)
}

func TestProductRepository_InsertProductLike_InsertErr(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	var userID uint64 = 1

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO user_favorite`).WithArgs(
		userID,
		prodTest.ID)
	mock.ExpectRollback()

	err = prodRepo.InsertProductLike(userID, prodTest.ID)
	assert.Error(t, err)
}

func TestProductRepository_InsertProductLike_UpdateErr(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	var userID uint64 = 1

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO user_favorite`).WithArgs(
		userID,
		prodTest.ID).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(`UPDATE product`).WithArgs(prodTest.ID)
	mock.ExpectRollback()

	err = prodRepo.InsertProductLike(userID, prodTest.ID)
	assert.Error(t, err)
}

func TestProductRepository_DeleteProductLike_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	var userID uint64 = 1

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE from user_favorite`).WithArgs(
		userID,
		prodTest.ID).WillReturnResult(driver.ResultNoRows)
	//mock.ExpectExec(`UPDATE product`).WithArgs(prodTest.ID).WillReturnResult(sqlmock.NewResult(int64(prodTest.ID), 1))
	mock.ExpectCommit()

	err = prodRepo.DeleteProductLike(userID, prodTest.ID)
	assert.NoError(t, err)
}

func TestProductRepository_DeleteProductLike_InsertErr(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	var userID uint64 = 1

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE from user_favorite`).WithArgs(
		userID,
		prodTest.ID)
	mock.ExpectRollback()

	err = prodRepo.DeleteProductLike(userID, prodTest.ID)
	assert.Error(t, err)
}

func TestProductRepository_DeleteProductLike_UpdateErr(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	var userID uint64 = 1

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE from user_favorite`).WithArgs(
		userID,
		prodTest.ID).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(`UPDATE product`).WithArgs(prodTest.ID)
	mock.ExpectRollback()

	err = prodRepo.DeleteProductLike(userID, prodTest.ID)
	assert.Error(t, err)
}

func TestProductRepository_UpdateTariff_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE product`).WithArgs(prodTest.Tariff, prodTest.ID).WillReturnResult(sqlmock.NewResult(int64(prodTest.ID), 1))
	mock.ExpectCommit()

	err = prodRepo.UpdateTariff(prodTest.ID, prodTest.Tariff)
	assert.NoError(t, err)
}

func TestProductRepository_UpdateTariff_Error(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE product`).WithArgs(prodTest.Tariff, prodTest.Address).WillReturnResult(sqlmock.NewResult(int64(prodTest.ID), 1))
	mock.ExpectCommit()

	err = prodRepo.UpdateTariff(prodTest.ID, prodTest.Tariff)
	assert.Error(t, err)
}

//func TestProductRepository_InsertPhoto_Success(t *testing.T) {
//	t.Parallel()
//
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer db.Close()
//
//	prodRepo := NewProductRepository(db)
//
//	mock.ExpectBegin()
//	mock.ExpectExec(`DELETE FROM product_images`).WithArgs(prodTest.ID).WillReturnResult(sqlmock.NewResult(int64(prodTest.ID), 1))
//	mock.ExpectExec(`INSERT INTO product_images`).WithArgs(
//		prodTest.ID,
//		prodTest.LinkImages[0]).WillReturnResult(driver.ResultNoRows)
//	mock.ExpectCommit()
//
//	err = prodRepo.InsertPhoto(prodTest)
//	assert.NoError(t, err)
//}

func TestProductRepository_InsertPhoto_DeleteErr(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM product_images`).WithArgs(prodTest.Address).WillReturnResult(sqlmock.NewResult(int64(prodTest.ID), 1))
	mock.ExpectRollback()

	err = prodRepo.InsertPhoto(prodTest)
	assert.Error(t, err)
}

func TestProductRepository_InsertPhoto_InserErr(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM product_images`).WithArgs(prodTest.ID).WillReturnResult(sqlmock.NewResult(int64(prodTest.ID), 1))
	mock.ExpectExec(`INSERT INTO product_images`).WithArgs(
		prodTest.Address,
		prodTest.LinkImages[0]).WillReturnResult(driver.ResultNoRows)
	mock.ExpectRollback()

	err = prodRepo.InsertPhoto(prodTest)
	assert.Error(t, err)
}

func TestProductRepository_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	answer := sqlmock.NewRows([]string{"id"}).AddRow(prodTest.ID)
	mock.ExpectBegin()
	mock.ExpectQuery(`UPDATE product`).WithArgs(prodTest.Name,
		prodTest.Amount,
		prodTest.Description,
		prodTest.Category,
		prodTest.Longitude,
		prodTest.Latitude,
		prodTest.Address,
		prodTest.ID).WillReturnRows(answer)
	mock.ExpectExec(`DELETE FROM product_images`).WithArgs(prodTest.ID).WillReturnResult(driver.RowsAffected(1))
	mock.ExpectExec(`INSERT INTO product_images`).WithArgs(prodTest.ID, sqlmock.AnyArg()).WillReturnResult(driver.RowsAffected(1))
	mock.ExpectCommit()

	err = prodRepo.Update(prodTest)

	assert.NoError(t, err)

	//error
	mock.ExpectBegin()
	mock.ExpectQuery(`UPDATE product`).WithArgs(prodTest.Name,
		prodTest.Amount,
		prodTest.Description,
		prodTest.Category,
		prodTest.Longitude,
		prodTest.Latitude,
		prodTest.Address,
		prodTest.ID).WillReturnError(sql.ErrConnDone)

	err = prodRepo.Update(prodTest)
	assert.Error(t, err)

	//error
	mock.ExpectBegin()
	mock.ExpectQuery(`UPDATE product`).WithArgs(prodTest.Name,
		prodTest.Amount,
		prodTest.Description,
		prodTest.Category,
		prodTest.Longitude,
		prodTest.Latitude,
		prodTest.Address,
		prodTest.ID).WillReturnRows(answer)
	mock.ExpectExec(`DELETE FROM product_images`).WithArgs(prodTest.ID).WillReturnError(sql.ErrConnDone)

	err = prodRepo.Update(prodTest)
	assert.Error(t, err)

	//error
	mock.ExpectBegin()
	mock.ExpectQuery(`UPDATE product`).WithArgs(prodTest.Name,
		prodTest.Amount,
		prodTest.Description,
		prodTest.Category,
		prodTest.Longitude,
		prodTest.Latitude,
		prodTest.Address,
		prodTest.ID).WillReturnRows(answer)
	mock.ExpectExec(`DELETE FROM product_images`).WithArgs(prodTest.ID).WillReturnResult(driver.RowsAffected(1))
	mock.ExpectExec(`INSERT INTO product_images`).WithArgs(prodTest.ID, sqlmock.AnyArg()).WillReturnError(sql.ErrConnDone)

	err = prodRepo.Update(prodTest)

	assert.Error(t, err)
}

func TestProductRepository_Close(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`Update product`).WithArgs(prodTest.ID).WillReturnResult(driver.RowsAffected(1))
	mock.ExpectCommit()

	err = prodRepo.Close(prodTest)
	assert.NoError(t, err)

	//error
	mock.ExpectBegin()
	mock.ExpectExec(`Update product`).WithArgs(prodTest.ID).WillReturnError(sql.ErrConnDone)

	err = prodRepo.Close(prodTest)
	assert.Error(t, err)
}

func TestProductRepository_SelectTrands(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	layout := "2006-01-02"
	date, _ := time.Parse(layout, prodTest.Date)

	linkStr := "{" + strings.Join(prodTest.LinkImages, ",") + "}"

	rows := sqlmock.NewRows([]string{"p.id", "p.name", "p.date", "p.amount", "array_agg(pi.img_link)", "uf.user_id", "p.tariff"})
	rows.AddRow(
		&prodTest.ID,
		&prodTest.Name,
		&date,
		&prodTest.Amount,
		&linkStr,
		&prodTest.OwnerID,
		&prodTest.Tariff)
	mock.ExpectQuery(`SELECT`).WithArgs(prodTest.OwnerID, uint64(1), uint64(2)).WillReturnRows(rows)

	_, err = prodRepo.SelectTrands([]uint64{1, 2}, &prodTest.OwnerID)
	assert.NoError(t, err)
}

func TestProductRepository_UpdateProductLikes(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE product`).WithArgs(1, prodTest.ID).WillReturnResult(driver.RowsAffected(1))
	mock.ExpectCommit()

	err = prodRepo.UpdateProductLikes(prodTest.ID, 1)
	assert.NoError(t, err)

	//error
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE product`).WithArgs(1, prodTest.ID).WillReturnError(sql.ErrConnDone)

	err = prodRepo.UpdateProductLikes(prodTest.ID, 1)
	assert.Error(t, err)
}

func TestProductRepository_UpdateProductViews(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE product`).WithArgs(1, prodTest.ID).WillReturnResult(driver.RowsAffected(1))
	mock.ExpectCommit()

	err = prodRepo.UpdateProductViews(prodTest.ID, 1)
	assert.NoError(t, err)

	//error
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE product`).WithArgs(1, prodTest.ID).WillReturnError(sql.ErrConnDone)

	err = prodRepo.UpdateProductViews(prodTest.ID, 1)
	assert.Error(t, err)
}

func TestProductRepository_SelectProductReviewers(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	answer := sqlmock.NewRows([]string{"u.id", "u.name", "u.avatar"}).AddRow(prodTest.OwnerID, prodTest.OwnerName, prodTest.OwnerLinkImages)
	mock.ExpectQuery(`SELECT`).WithArgs(prodTest.ID, prodTest.OwnerID).WillReturnRows(answer)

	_, err = prodRepo.SelectProductReviewers(prodTest.ID, prodTest.OwnerID)
	assert.NoError(t, err)

	//error
	mock.ExpectQuery(`SELECT`).WithArgs(prodTest.ID, prodTest.OwnerID).WillReturnError(sql.ErrConnDone)

	_, err = prodRepo.SelectProductReviewers(prodTest.ID, prodTest.OwnerID)
	assert.Error(t, err)
}

func TestProductRepository_InsertProductBuyer(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE product`).WithArgs(prodTest.ID, prodTest.OwnerID).WillReturnResult(driver.RowsAffected(1))
	mock.ExpectCommit()

	err = prodRepo.InsertProductBuyer(prodTest.ID, prodTest.OwnerID)
	assert.NoError(t, err)

	//error
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE product`).WithArgs(prodTest.ID, prodTest.OwnerID).WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err = prodRepo.InsertProductBuyer(prodTest.ID, prodTest.OwnerID)
	assert.Error(t, err)
}

func TestProductRepository_CheckProductReview(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodRepo := NewProductRepository(db)

	answer := sqlmock.NewRows([]string{"owner_id", "seller_left_review"}).AddRow(prodTest.OwnerID, false)
	mock.ExpectQuery(`SELECT`).WithArgs(prodTest.ID).WillReturnRows(answer)

	_, err = prodRepo.CheckProductReview(prodTest.ID, "seller", prodTest.OwnerID)
	assert.NoError(t, err)

	//error
	mock.ExpectQuery(`SELECT`).WithArgs(prodTest.ID).WillReturnError(sql.ErrConnDone)

	_, err = prodRepo.CheckProductReview(prodTest.ID, "buyer", prodTest.OwnerID)
	assert.Error(t, err)
}