package usecase

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/favorite"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/jackc/pgx"
)

type FavoriteUsecase struct {
	favoriteRepo favorite.FavoriteRepository
}

func NewFavoriteUsecase(repo favorite.FavoriteRepository) favorite.FavoriteUsecase {
	return &FavoriteUsecase{
		favoriteRepo: repo,
	}
}

func (pu *FavoriteUsecase) LikeProduct(userID uint64, productID uint64) *errors.Error {
	err := pu.favoriteRepo.InsertProduct(userID, productID)

	if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
		return errors.Cause(errors.ProductAlreadyLiked)
	}

	if err != nil {
		return errors.UnexpectedInternal(err)
	}

	return nil
}

func (pu *FavoriteUsecase) DislikeProduct(userID uint64, productID uint64) *errors.Error {
	err := pu.favoriteRepo.DeleteProduct(userID, productID)
	if err != nil {
		return errors.UnexpectedInternal(err)
	}

	return nil
}

func (pu *FavoriteUsecase) GetUserFavorite(userID uint64, content *models.Page) ([]*models.ProductListData, *errors.Error) {
	products, err := pu.favoriteRepo.SelectUserFavorite(userID, content)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	if len(products) == 0 {
		return []*models.ProductListData{}, nil
	}

	return products, nil
}
