package repository

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

var prodTest = &models.ProductData{
	ID:           0,
	Name:         "tovar",
	Date:         "2010-02-02",
	Amount:       10000,
	Description:  "Description product aaaaa",
	Category:     "0",
	OwnerID:      0,
	Longitude: 	  1,
	Latitude:     1,
	Address:      "Address",
	LinkImages:   []string{"test_str.jpg"},
	Tariff: 0,
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
		"u.surname", "u.avatar", "p.likes", "p.views", "p.longitude", "p.latitude", "p.address", "array_agg(pi.img_link)", "p.tariff"})
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
		&prodTest.Likes,
		&prodTest.Views,
		&prodTest.Longitude,
		&prodTest.Latitude,
		&prodTest.Address,
		&linkStr,
		&prodTest.Tariff)
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT`).WithArgs(prodTest.ID).WillReturnRows(rows)
	mock.ExpectExec(`UPDATE product`).WithArgs(prodTest.ID).WillReturnResult(sqlmock.NewResult(int64(prodTest.ID), 1))
	mock.ExpectCommit()

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
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT`).WithArgs(prodTest.ID).WillReturnRows(rows)
	mock.ExpectRollback()

	_, err = prodRepo.SelectByID(prodTest.ID)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductRepository_SelectByID_UpdateErr(t *testing.T) {
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
		"u.surname", "u.avatar", "p.likes", "p.views", "p.longitude", "p.latitude", "p.address", "array_agg(pi.img_link)", "p.tariff"})
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
		&prodTest.Likes,
		&prodTest.Views,
		&prodTest.Longitude,
		&prodTest.Latitude,
		&prodTest.Address,
		&linkStr,
		&prodTest.Tariff)
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT`).WithArgs(prodTest.ID).WillReturnRows(rows)
	mock.ExpectExec(`UPDATE product`).WithArgs(prodTest.ID)
	mock.ExpectRollback()

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
	mock.ExpectExec(`UPDATE product`).WithArgs(prodTest.ID).WillReturnResult(sqlmock.NewResult(int64(prodTest.ID), 1))
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
	mock.ExpectExec(`UPDATE product`).WithArgs(prodTest.ID).WillReturnResult(sqlmock.NewResult(int64(prodTest.ID), 1))
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

func TestProductRepository_InsertPhoto_Success(t *testing.T) {
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
		prodTest.ID,
		prodTest.LinkImages[0]).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = prodRepo.InsertPhoto(prodTest)
	assert.NoError(t, err)
}

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