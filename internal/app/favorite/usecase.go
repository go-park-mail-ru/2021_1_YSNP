package favorite

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

type FavoriteUsecase interface {
	LikeProduct(userID uint64, productID uint64) *errors.Error
	DislikeProduct(userID uint64, productID uint64) *errors.Error
	GetUserFavorite(userID uint64, content *models.Page) ([]*models.ProductListData, *errors.Error)
}
