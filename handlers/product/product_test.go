package product

import (
	_tmpDB "2021_1_YSNP/tmp_database"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProductIDHandler_ProductIDHandlerSuccess(t *testing.T) {
	_tmpDB.InitDB()

	expectedJSON := `{"id":0,"name":"iphone","date":"2000-10-10","amount":12000,"linkImages":["http://89.208.199.170:8080/static/product/pic4.jpeg","http://89.208.199.170:8080/static/product/pic5.jpeg","http://89.208.199.170:8080/static/product/pic6.jpeg"],"description":"eto iphone","ownerId":0,"ownerName":"Sergey","ownerSurname":"Alehin","views":1000,"likes":20}`

	r := httptest.NewRequest("GET", "/api/v1/product/0", nil)
	w := httptest.NewRecorder()

	ProductIDHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestProductIDHandler_ProductIDHandlerNoProduct(t *testing.T) {
	expectedJSON := `{"message":"No product with this id."}`

	r := httptest.NewRequest("GET", "/api/v1/product/3", nil)
	w := httptest.NewRecorder()

	ProductIDHandler(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("status is not 404")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestProductCreateHandler_ProductCreateHandlerNotAuth(t *testing.T) {
	_tmpDB.InitDB()

	var expectedJSON = `{"message":"User not authorised or not found"}`


	r := httptest.NewRequest("POST", "/api/v1/product/create", nil)
	w := httptest.NewRecorder()

	ProductCreateHandler(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Error("status is not 401")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestProductCreateHandler_ProductCreateHandlerWrongRequest(t *testing.T) {
	_tmpDB.InitDB()

	var expectedJSON = `{"message":"invalid character '}' looking for beginning of object key string"}`

	var byteData = bytes.NewReader([]byte(`{
			"telephone" : "+79990009900",
			"password" : "Qwerty12",
		}`))

	r := httptest.NewRequest("POST", "/api/v1/product/create", byteData)
	r.AddCookie(&http.Cookie{Name:"session_id", Value: _tmpDB.NewSession("+79990009900")})
	w := httptest.NewRecorder()

	ProductCreateHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("status is not 400")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

//func TestProductCreateHandler_ProductCreateHandlerSuccess(t *testing.T) {
//	_tmpDB.InitDB()
//
//	var expectedJSON = `{"message":"invalid character '}' looking for beginning of object key string"}`
//
//	var byteData = bytes.NewReader([]byte(`{
//			"telephone" : "+79990009900",
//			"password" : "Qwerty12",
//		}`))
//
//	r := httptest.NewRequest("POST", "/api/v1/product/create", byteData)
//	r.AddCookie(&http.Cookie{Name:"session_id", Value: _tmpDB.NewSession("+79990009900")})
//	w := httptest.NewRecorder()
//
//	ProductCreateHandler(w, r)
//
//	if w.Code != http.StatusOK {
//		t.Error("status is not ok")
//	}
//
//	bytes, _ := ioutil.ReadAll(w.Body)
//	if strings.Trim(string(bytes), "\n") != expectedJSON {
//		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
//	}
//}