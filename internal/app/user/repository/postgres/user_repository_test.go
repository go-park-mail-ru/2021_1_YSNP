package repository

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

var userTest = &models.UserData{
	ID:         0,
	Name:       "Максим",
	Surname:    "Торжков",
	Sex:        "male",
	Email:      "a@a.ru",
	Telephone:  "+79169230768",
	Password:   "Qwerty12",
	DateBirth:  "2021-03-08",
	Latitude:   0,
	Longitude:  0,
	Radius:     0,
	Address:    "",
	LinkImages: "",
	Rating: 0,
}

var userOauthText = &models.UserOAuthRequest{
	ID: 0,
	FirstName: "Максим",
	LastName: "Торжков",
	Photo: "",
	UserOAuthID: 34,
	UserOAuthType: "",
}

func TestUserRepository_SelectByID_OK(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)

	layout := "2006-01-02"
	time, _ := time.Parse(layout, userTest.DateBirth)

	rows := sqlmock.NewRows([]string{"id", "email", "telephone", "password", "name", "surname", "sex", "birthdate", "latitude", "longitude", "radius", "address", "avatar", "score", "reviews"})
	rows.AddRow(
		userTest.ID,
		userTest.Email,
		userTest.Telephone,
		userTest.Password,
		userTest.Name,
		userTest.Surname,
		userTest.Sex,
		time,
		userTest.Latitude,
		userTest.Longitude,
		userTest.Radius,
		userTest.Address,
		userTest.LinkImages,
		userTest.Rating,
		0)
	mock.ExpectQuery(`SELECT`).WithArgs(userTest.ID).WillReturnRows(rows)

	user, err := userRepo.SelectByID(userTest.ID)
	assert.Equal(t, userTest, user)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_SelectByID_Error(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "email", "telephone", "password", "name", "surname", "sex", "birthdate", "latitude", "longitude", "radius", "address", "avatar"})
	rows.AddRow(
		userTest.ID,
		userTest.Email,
		userTest.Telephone,
		userTest.Password,
		userTest.Name,
		userTest.Surname,
		userTest.Sex,
		userTest.DateBirth,
		userTest.Latitude,
		userTest.Longitude,
		userTest.Radius,
		userTest.Address,
		userTest.LinkImages)
	mock.ExpectQuery(`SELECT`).WithArgs(userTest.ID).WillReturnRows(rows)

	_, err = userRepo.SelectByID(userTest.ID)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_SelectByTelephone_OK(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)

	layout := "2006-01-02"
	time, _ := time.Parse(layout, userTest.DateBirth)

	rows := sqlmock.NewRows([]string{"id", "email", "telephone", "password", "name", "surname", "sex", "birthdate", "latitude", "longitude", "radius", "address", "avatar", "score", "reviews"})
	rows.AddRow(
		userTest.ID,
		userTest.Email,
		userTest.Telephone,
		userTest.Password,
		userTest.Name,
		userTest.Surname,
		userTest.Sex,
		time,
		userTest.Latitude,
		userTest.Longitude,
		userTest.Radius,
		userTest.Address,
		userTest.LinkImages,
		userTest.Rating,
		0)
	mock.ExpectQuery(`SELECT`).WithArgs(userTest.Telephone).WillReturnRows(rows)

	user, err := userRepo.SelectByTelephone(userTest.Telephone)
	assert.Equal(t, userTest, user)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_SelectByTelephone_Error(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "email", "telephone", "password", "name", "surname", "sex", "birthdate", "latitude", "longitude", "radius", "address", "avatar", })
	rows.AddRow(
		userTest.ID,
		userTest.Email,
		userTest.Telephone,
		userTest.Password,
		userTest.Name,
		userTest.Surname,
		userTest.Sex,
		userTest.DateBirth, // wrong date type (time.Time -> string)
		userTest.Latitude,
		userTest.Longitude,
		userTest.Radius,
		userTest.Address,
		userTest.LinkImages)
	mock.ExpectQuery(`SELECT`).WithArgs(userTest.Telephone).WillReturnRows(rows)

	_, err = userRepo.SelectByTelephone(userTest.Telephone)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Insert_OK(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)

	mock.ExpectBegin()
	answer := sqlmock.NewRows([]string{"id"}).AddRow(userTest.ID)
	mock.ExpectQuery(`INSERT INTO users`).WithArgs(
		userTest.Email,
		userTest.Telephone,
		userTest.Password,
		userTest.Name,
		userTest.Surname,
		userTest.Sex,
		userTest.DateBirth,
		userTest.LinkImages).WillReturnRows(answer)
	mock.ExpectCommit()

	err = userRepo.Insert(userTest)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Insert_Error(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)

	mock.ExpectBegin()
	answer := sqlmock.NewRows([]string{"id"}).AddRow(userTest.Telephone) //scan wrong type
	mock.ExpectQuery(`INSERT INTO users`).WithArgs(
		userTest.Email,
		userTest.Telephone,
		userTest.Password,
		userTest.Name,
		userTest.Surname,
		userTest.Sex,
		userTest.DateBirth,
		userTest.LinkImages).WillReturnRows(answer)
	mock.ExpectRollback()

	err = userRepo.Insert(userTest)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Update_OK(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE users`).WithArgs(
		userTest.ID,
		userTest.Email,
		userTest.Telephone,
		userTest.Password,
		userTest.Name,
		userTest.Surname,
		userTest.Sex,
		userTest.DateBirth,
		userTest.Latitude,
		userTest.Longitude,
		userTest.Radius,
		userTest.Address,
		userTest.LinkImages).WillReturnResult(sqlmock.NewResult(int64(userTest.ID), 1))
	mock.ExpectCommit()

	err = userRepo.Update(userTest)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Update_Error(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE users`).WithArgs(
		userTest.ID,
		userTest.Email,
		userTest.Telephone,
		userTest.Password,
		userTest.Name,
		userTest.Surname,
		userTest.Sex,
		userTest.DateBirth,
		userTest.Latitude,
		userTest.Longitude,
		userTest.Radius,
		userTest.Address,
		userTest.LinkImages)
	mock.ExpectRollback()

	err = userRepo.Update(userTest)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_InsertOAuth_Error(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)

	mock.ExpectBegin()
	answer := sqlmock.NewRows([]string{"id"}).AddRow(userTest.ID) //scan wrong type
	mock.ExpectQuery(`INSERT INTO users`).WithArgs(
		userOauthText.FirstName,
		userOauthText.LastName,
		userOauthText.Photo).WillReturnRows(answer)
	mock.ExpectRollback()

	err = userRepo.InsertOAuth(userOauthText)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_InsertOAuth_OK(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)

	mock.ExpectBegin()
	answer := sqlmock.NewRows([]string{"id"}).AddRow(userOauthText.ID)
	mock.ExpectQuery(`INSERT INTO users`).WithArgs(
		userOauthText.FirstName,
		userOauthText.LastName,
		userOauthText.Photo).WillReturnRows(answer)

	mock.ExpectExec(`INSERT INTO users_oauth`).WithArgs(
		userOauthText.ID,
		userOauthText.UserOAuthType,
		userOauthText.UserOAuthID).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = userRepo.InsertOAuth(userOauthText)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}


func TestUserRepository_SelectByOAuthID__Error(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "last_name", "first_name", "photo_max", "user_oauth_id", "user_oauth_type" })
	rows.AddRow(
		userOauthText.ID,
		userOauthText.FirstName,
		userOauthText.LastName,
		userOauthText.Photo,
		userOauthText.UserOAuthID,
		userOauthText.UserOAuthType)
	mock.ExpectQuery(`SELECT`).WithArgs(userOauthText.UserOAuthID).WillReturnRows(rows)

	id := userRepo.SelectByOAuthID(userOauthText.UserOAuthID)
	assert.Equal(t, uint64(0), id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_SelectByOAuthID_OK(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)


	rows := sqlmock.NewRows([]string{"user_id" })
	rows.AddRow(
		userOauthText.ID)
	mock.ExpectQuery(`SELECT`).WithArgs(userOauthText.UserOAuthID).WillReturnRows(rows)

	user := userRepo.SelectByOAuthID(userOauthText.UserOAuthID)
	assert.Equal(t, userOauthText.ID, user)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}