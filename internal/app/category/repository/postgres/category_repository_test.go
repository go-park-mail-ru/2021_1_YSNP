package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCategoryRepository_SelectCategories(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	catRepo := NewCategoryRepository(db)

	testTitle := "title"
	answer := sqlmock.NewRows([]string{"title"}).AddRow(testTitle)
	mock.ExpectQuery(`SELECT`).WillReturnRows(answer)

	_, err = catRepo.SelectCategories()
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//error
	mock.ExpectQuery(`SELECT`).WillReturnError(sql.ErrNoRows)

	_, err = catRepo.SelectCategories()
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
