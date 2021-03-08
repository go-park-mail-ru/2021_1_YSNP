package main_page

import (
	_tmpDB "2021_1_YSNP/tmp_database"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


func TestMainPageHandler_MainPageHandlerSuccess(t *testing.T) {
	_tmpDB.InitDB()

	var expectedJSON = `{"product_list":[{"id":0,"name":"iphone","date":"2000-10-10","amount":12000,"linkImages":["http://89.208.199.170:8080/static/product/pic4.jpeg","http://89.208.199.170:8080/static/product/pic5.jpeg","http://89.208.199.170:8080/static/product/pic6.jpeg"]},{"id":1,"name":"iphone 10","date":"2000-10-10","amount":12001,"linkImages":["http://89.208.199.170:8080/static/product/pic1.jpeg","http://89.208.199.170:8080/static/product/pic2.jpeg","http://89.208.199.170:8080/static/product/pic3.jpeg"]}]}`


	r := httptest.NewRequest("GET", "/api/v1/", nil)
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
