package repository

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

func TestSearchRepository_SelectByFilter(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prodTest := &models.ProductListData{
		ID:         0,
		Name:       "Product",
		Date:       "2010-02-02T15:04:05Z",
		Amount:     10000,
		LinkImages: []string{"test_png.png"},
		UserLiked:  true,
		Tariff:     0,
	}

	data := &models.Search{
		Category:   "",
		Date:       "",
		FromAmount: 0,
		ToAmount:   200,
		Radius:     1,
		Latitude:   0,
		Longitude:  0,
		Search:     "",
		Sorting:    "",
		From:       0,
		Count:      0,
	}

	var userID uint64 = 1

	layout := "2006-01-02T15:04:05Z"
	date, _ := time.Parse(layout, prodTest.Date)

	linkStr := "{" + strings.Join(prodTest.LinkImages, ",") + "}"

	searchRep := NewSearchRepository(db)

	rows := sqlmock.NewRows([]string{"p.id", "p.name", "p.date", "p.amount", "array_agg(pi.img_link)", "uf.user_id", "p.tariff"})
	rows.AddRow(
		&prodTest.ID,
		&prodTest.Name,
		&date,
		&prodTest.Amount,
		&linkStr,
		&userID,
		&prodTest.Tariff)
	mock.ExpectQuery(`SELECT`).WithArgs(userID,
		"%"+data.Search+"%",
		"%"+data.Category+"%",
		data.FromAmount,
		data.ToAmount,
		"SRID=4326; POINT("+fmt.Sprintf("%f", data.Longitude)+" "+fmt.Sprintf("%f", data.Latitude)+")",
		data.Radius,
		data.Count,
		data.From*data.Count).WillReturnRows(rows)

	prod, err := searchRep.SelectByFilter(&userID, data)
	assert.Equal(t, prodTest, prod[0])
	assert.NoError(t, err)
}
