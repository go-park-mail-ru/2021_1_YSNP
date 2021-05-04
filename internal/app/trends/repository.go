package trends

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

//go:generate mockgen -destination=./mocks/mock_session_repo.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session SessionRepository

type TrendsRepository interface {
	InsertOrUpdate(ui *models.Trends) error
	CreateTrendsProducts(userID uint64) error
	GetTrendsProducts(userID uint64) ([]uint64, error)
}

