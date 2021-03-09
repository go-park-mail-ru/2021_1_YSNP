package main_page

import (
	_tmpDB "2021_1_YSNP/tmp_database"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestMainPageHandler_MainPageHandlerSuccess(t *testing.T) {
	_tmpDB.InitDB()

	r := httptest.NewRequest("GET", "/api/v1/", nil)
	w := httptest.NewRecorder()

	MainPageHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}
}
