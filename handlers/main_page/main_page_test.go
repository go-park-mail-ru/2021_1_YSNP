package main_page

import (
	_tmpDB "2021_1_YSNP/tmp_database"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var expectedJSON = `{"product_list":[{"id":0,"name":"iphone","date":"2000-10-10","amount":12000,"link_images":null},{"id":1,"name":"iphone 10","date":"2000-10-10","amount":12001,"link_images":null}]}`

func TestMainPageHandler_MainPageHandlerSuccess(t *testing.T) {
	_tmpDB.InitDB()

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	MainPageHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}
