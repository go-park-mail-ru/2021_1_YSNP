package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

type UserRepository struct {
	dbConn *sql.DB
}

func NewUserRepository(conn *sql.DB) user.UserRepository{
	return &UserRepository{
			dbConn: conn,
	}
}

func (ur *UserRepository) Insert(user *models.UserData) error {
	panic("implement me")
}

func (ur *UserRepository) SelectByTelephone(telephone string) (*models.UserData, error) {
	panic("implement me")
}

func (ur *UserRepository) SelectByID(userID uint64) (*models.UserData, error) {
	panic("implement me")
}

func (ur *UserRepository) Update(user *models.UserData) error {
	panic("implement me")
}