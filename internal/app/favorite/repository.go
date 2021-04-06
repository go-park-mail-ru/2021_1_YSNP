package favorite

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

type FavoriteRepository interface {
	InsertProduct(userID uint64, productID uint64) error
	DeleteProduct(userID uint64, productID uint64) error
	SelectUserFavorite(userID uint64, content *models.Page) ([]*models.ProductListData, error)
}
