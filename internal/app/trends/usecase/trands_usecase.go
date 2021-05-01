package usecase

import (
	"strings"
	"github.com/bbalet/stopwords"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends"
	errors2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

type TrendsUsecase struct {
	TrendsRepo trends.TrendsRepository
}

func NewTrendsUsecase(repo trends.TrendsRepository) trends.TrendsUsecase {
	return &TrendsUsecase{
		TrendsRepo: repo,
	}
}


func (tu *TrendsUsecase) InsertOrUpdate(ui *models.UserInterested) *errors2.Error {
	cleanContent := stopwords.CleanString(ui.Text, "ru", true)
	sn := strings.TrimSpace(cleanContent)
	s := strings.FieldsFunc(sn, func(r rune) bool { return strings.ContainsRune(" .,:-", r) })

	ua := &models.Trends{}
	ua.UserID = ui.UserID
	for _, item := range s {
		ua.Popular = append(ua.Popular, models.Popular{
			Title: item,
			Count: 1,
		})
	}
	err := tu.TrendsRepo.InsertOrUpdate(ua)
	if err != nil {
		return errors2.UnexpectedInternal(err)
	}
	err = tu.TrendsRepo.CreateTrendsProducts(ui.UserID)
	if err != nil {
		return errors2.UnexpectedInternal(err)
	}

	return nil
}


