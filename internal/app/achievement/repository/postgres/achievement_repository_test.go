package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAchievementRepository_GetUserAchievements(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	achRepo := NewAchievementRepository(db)

	testTitle := "title"
	testDescription := "description"
	testDate := "date"
	testLinkPic := "link_pic"
	answer := sqlmock.NewRows([]string{"title", "description", "date", "link_pic"}).AddRow(testTitle, testDescription, testDate, testLinkPic)
	mock.ExpectQuery(`SELECT`).WillReturnRows(answer)

	_, err = achRepo.GetUserAchievements(0)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//error
	mock.ExpectQuery(`SELECT`).WillReturnError(sql.ErrNoRows)

	_, err = achRepo.GetUserAchievements(0)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
