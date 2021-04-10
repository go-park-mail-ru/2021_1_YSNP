package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
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
		user.Email,
		user.Telephone,
		user.Password,
		user.Name,
		user.Surname,
		user.Sex,
		user.DateBirth,
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
				avatar
				FROM users
				WHERE telephone=$1`,
		telephone)

	var date time.Time

	err := query.Scan(
		&user.ID,
		&user.Email,
		&user.Telephone,
		&user.Password,
		&user.Name,
		&user.Surname,
		&user.Sex,
		&date,
		&user.Latitude,
		&user.Longitude,
		&user.Radius,
		&user.Address,
		&user.LinkImages)
	if err != nil {
		return nil, err
	}

	user.DateBirth = date.Format("2006-01-02")

	return user, nil
}

func (ur *UserRepository) SelectByID(userID uint64) (*models.UserData, error) {
	user := &models.UserData{}

	query := ur.dbConn.QueryRow(
		`
				SELECT id, email, telephone, password, 
			    name, surname, sex, birthdate, 
				latitude, longitude, radius, address, avatar
				FROM users
				WHERE id=$1`,
		userID)

	var date time.Time

	err := query.Scan(
		&user.ID,
		&user.Email,
		&user.Telephone,
		&user.Password,
		&user.Name,
		&user.Surname,
		&user.Sex,
		&date,
		&user.Latitude,
		&user.Longitude,
		&user.Radius,
		&user.Address,
		&user.LinkImages)
	if err != nil {
		return nil, err
	}

	user.DateBirth = date.Format("2006-01-02")

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
		user.Email,
		user.Password,
		user.Name,
		user.Surname,
		user.Sex,
		user.DateBirth,
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
