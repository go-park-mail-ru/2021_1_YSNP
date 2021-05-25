package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/null_value"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
)

type UserRepository struct {
	dbConn *sql.DB
}

func NewUserRepository(conn *sql.DB) user.UserRepository {
	return &UserRepository{
		dbConn: conn,
	}
}

func (ur *UserRepository) Insert(user *models.UserData) error {
	tx, err := ur.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	query := tx.QueryRow(
		`
				INSERT INTO users(email, telephone, password, name, surname, sex, birthdate, avatar)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING id`,
		null_value.NewNullString(user.Email),
		null_value.NewNullString(user.Telephone),
		null_value.NewNullString(user.Password),
		user.Name,
		user.Surname,
		null_value.NewNullString(user.Sex),
		null_value.NewNullString(user.DateBirth),
		user.LinkImages)

	err = query.Scan(&user.ID)
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

func (ur *UserRepository) SelectByTelephone(telephone string) (*models.UserData, error) {
	user := &models.UserData{}

	query := ur.dbConn.QueryRow(
		`
				SELECT id, email, telephone, password, 
				name, surname, sex, birthdate, 
				latitude, longitude, radius, address,
				avatar, score, reviews
				FROM users
				WHERE telephone=$1`,
		telephone)

	var nullEmail sql.NullString
	var nullTelephone sql.NullString
	var nullPassword sql.NullString
	var nullSex sql.NullString
	var nullDate sql.NullTime
	var reviews int

	err := query.Scan(
		&user.ID,
		&nullEmail,
		&nullTelephone,
		&nullPassword,
		&user.Name,
		&user.Surname,
		&nullSex,
		&nullDate,
		&user.Latitude,
		&user.Longitude,
		&user.Radius,
		&user.Address,
		&user.LinkImages,
		&user.Rating,
		&reviews)
	if err != nil {
		return nil, err
	}

	user.Email = null_value.NewStringFromNull(nullEmail)
	user.Telephone = null_value.NewStringFromNull(nullTelephone)
	user.Password = null_value.NewStringFromNull(nullPassword)
	user.Sex = null_value.NewStringFromNull(nullSex)
	user.Rating = user.Rating / float64(reviews)

	if nullDate.Valid {
		date := nullDate.Time
		user.DateBirth = date.Format("2006-01-02")
	}

	return user, nil
}

func (ur *UserRepository) SelectByID(userID uint64) (*models.UserData, error) {
	user := &models.UserData{}

	query := ur.dbConn.QueryRow(
		`
				SELECT id, email, telephone, password, 
			    name, surname, sex, birthdate, 
				latitude, longitude, radius, address, avatar, score, reviews
				FROM users
				WHERE id=$1`,
		userID)

	var nullEmail sql.NullString
	var nullTelephone sql.NullString
	var nullPassword sql.NullString
	var nullSex sql.NullString
	var nullDate sql.NullTime
	var reviews int

	err := query.Scan(
		&user.ID,
		&nullEmail,
		&nullTelephone,
		&nullPassword,
		&user.Name,
		&user.Surname,
		&nullSex,
		&nullDate,
		&user.Latitude,
		&user.Longitude,
		&user.Radius,
		&user.Address,
		&user.LinkImages,
		&user.Rating,
		&reviews)
	if err != nil {
		return nil, err
	}

	user.Email = null_value.NewStringFromNull(nullEmail)
	user.Telephone = null_value.NewStringFromNull(nullTelephone)
	user.Password = null_value.NewStringFromNull(nullPassword)
	user.Sex = null_value.NewStringFromNull(nullSex)
	user.Rating = user.Rating / float64(reviews)

	if nullDate.Valid {
		date := nullDate.Time
		user.DateBirth = date.Format("2006-01-02")
	}

	return user, nil
}

func (ur *UserRepository) Update(user *models.UserData) error {
	tx, err := ur.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`
				UPDATE users
				SET email = $2, telephone = $3, password = $4, 
				name = $5, surname = $6, sex = $7, birthdate = $8,
				latitude = $9, longitude = $10, radius = $11, address = $12,
				avatar = $13
				WHERE id = $1;`,
		user.ID,
		null_value.NewNullString(user.Email),
		null_value.NewNullString(user.Telephone),
		null_value.NewNullString(user.Password),
		user.Name,
		user.Surname,
		null_value.NewNullString(user.Sex),
		null_value.NewNullString(user.DateBirth),
		user.Latitude,
		user.Longitude,
		user.Radius,
		user.Address,
		user.LinkImages)

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

func (ur *UserRepository) InsertOAuth(userOAuth *models.UserOAuthRequest) error {
	tx, err := ur.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	query := tx.QueryRow(
		`
				INSERT INTO users(name, surname, avatar)
				VALUES ($1, $2, $3)
				RETURNING id`,
		userOAuth.FirstName,
		userOAuth.LastName,
		userOAuth.Photo)

	err = query.Scan(&userOAuth.ID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	_, err = tx.Exec(
		`
				INSERT INTO users_oauth
                (user_id, oauth_type, oauth_id)
                VALUES ($1, $2, $3)`,
		userOAuth.ID,
		userOAuth.UserOAuthType,
		userOAuth.UserOAuthID)
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

func (ur *UserRepository) SelectByOAuthID(userOAuthID float64) uint64 {
	query := ur.dbConn.QueryRow(
		`
				SELECT user_id
				FROM users_oauth
				WHERE oauth_id =$1`,
		userOAuthID)

	var userID uint64

	err := query.Scan(
		&userID)
	if err != nil {
		return 0
	}

	return userID
}
