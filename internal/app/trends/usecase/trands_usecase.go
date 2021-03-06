package usecase

import (
	"strings"
	"time"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends"
)

type TrendsUsecase struct {
	TrendsRepo trends.TrendsRepository
}

func NewTrendsUsecase(repo trends.TrendsRepository) trends.TrendsUsecase {
	return &TrendsUsecase{
		TrendsRepo: repo,
	}
}

func checkSufix(word string) bool {
	if len(word) < 5 {
		return true
	}
	stop := []string{"ими", "ыми", "его", "ого", "ему", "ому", "ее", "ие",
		"ые", "ое", "ей", "ий", "ый", "ой", "ем", "им", "ым",
		"ом", "их", "ых", "ую", "юю", "ая", "яя", "ою", "ею"}

	for _, item := range stop {
		suf := word[len(word)-len(item):]
		if suf == item {
			return false
		}
	}
	return true
}

func (tu *TrendsUsecase) InsertOrUpdate(ui *models.UserInterested) *errors.Error {
	cleanContent := stopwords.CleanString(ui.Text, "ru", true)
	sn := strings.TrimSpace(cleanContent)
	s := strings.FieldsFunc(sn, func(r rune) bool { return strings.ContainsRune(" .,:-", r) })

	ua := &models.Trends{}
	ua.UserID = ui.UserID
	for _, item := range s {
		if !checkSufix(item) {
			continue
		}

		stemmed, err := snowball.Stem(item, "russian", true)
		if err == nil {
			ua.Popular = append(ua.Popular, models.Popular{
				Title: stemmed,
				Count: 1,
				Date:  time.Now(),
			})
		}
	}
	err := tu.TrendsRepo.InsertOrUpdate(ua)
	if err != nil {
		return errors.UnexpectedInternal(err)
	}
	//nolint:errcheck
	go tu.TrendsRepo.CreateTrendsProducts(ui.UserID)
	return nil
}
