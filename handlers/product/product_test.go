package product

import (
	_tmpDB "2021_1_YSNP/tmp_database"
	"bytes"
	"github.com/gorilla/mux"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestProductIDHandler_ProductIDHandlerSuccess(t *testing.T) {
	_tmpDB.InitDB()

	expectedJSON := `{"id":0,"name":"iphone","date":"2000-10-10","amount":5994,"linkImages":["http://89.208.199.170:8080/static/product/pic10.jpeg","http://89.208.199.170:8080/static/product/pic7.jpeg","http://89.208.199.170:8080/static/product/pic3.jpeg"],"description":"Ясность нашей позиции очевидна: перспективное планирование играет определяющее значение для благоприятных перспектив. Противоположная точка зрения подразумевает, что сторонники тоталитаризма в науке неоднозначны и будут объективно рассмотрены соответствующими инстанциями.","category":"Автомобили","ownerId":0,"ownerName":"Sergey","ownerSurname":"Alehin","views":23,"likes":634}`

	r := httptest.NewRequest("GET", "/api/v1/product/0", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
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
	expectedJSON := `{"message":"no product with this id"}`

	r := httptest.NewRequest("GET", "/api/v1/product/30000000", nil)
	r.AddCookie(&http.Cookie{Name: "session_id", Value: _tmpDB.NewSession("+79990009900")})
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

	var expectedJSON = `{"message":"user not authorised or not found"}`

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
	r.AddCookie(&http.Cookie{Name: "session_id", Value: _tmpDB.NewSession("+79990009900")})
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

func TestProductCreateHandler_ProductCreateHandlerSuccess(t *testing.T) {
	_tmpDB.InitDB()

	var expectedJSON = `{"message":"Successful creation."}`

	var byteData = bytes.NewReader([]byte(`{
				"name":"ddsdsdsdsd",
				"description":"wefdsfdffsf",
				"amount":2323323,
				"linkImages":["http://89.208.199.170:8080/static/product/12fd0fce-6022-4e1c-9858-b64ff914f9cf.jpg"],
				"category":"Электроника"
	}`))

	r := httptest.NewRequest("POST", "/api/v1/product/create", byteData)
	r.AddCookie(&http.Cookie{Name: "session_id", Value: _tmpDB.NewSession("+79990009900")})
	w := httptest.NewRecorder()

	ProductCreateHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestUploadPhotoHandler_UploadPhotoHandlerWrongContentType(t *testing.T) {
	expectedJSON := `{"message":"request Content-Type isn't multipart/form-data"}`

	//var byteData = bytes.NewReader([]byte(`{"linkImages":"http://89.208.199.170:8080/static/avatar/b3c098f5-94d8-4bb9-8e56-bc626e60aab7.jpg"}`))

	r := httptest.NewRequest("POST", "/api/v1/product/upload", nil)
	r.AddCookie(&http.Cookie{Name: "session_id", Value: _tmpDB.NewSession("+79990009900")})
	w := httptest.NewRecorder()

	UploadPhotoHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestUploadPhotoHandler_UploadPhotoHandlerSucces(t *testing.T) {
	path := "../../static/avatar/test-avatar.jpg"
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("photos", path)
	if err != nil {
		t.Fatal(err)
	}
	sample, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	text, _ := ioutil.ReadAll(sample)
	part.Write(text)
	writer.Close()
	sample.Close()

	r := httptest.NewRequest("POST", "/api/v1/product/upload", body)

	r.Header.Add("Content-Type", writer.FormDataContentType())
	r.AddCookie(&http.Cookie{Name: "session_id", Value: _tmpDB.NewSession("+79990009900")})
	w := httptest.NewRecorder()

	UploadPhotoHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
	}
}

func TestUploadPhotoHandler_UploadPhotoHandlerNoFile(t *testing.T) {
	expectedJSON := `{"message":"http: no such file"}`

	path := "../../static/avatar/test-avatar.jpg"
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file_", path)
	if err != nil {
		t.Error(err)
	}
	sample, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}
	text, err := ioutil.ReadAll(sample)
	if err != nil {
		t.Error(err)
	}
	part.Write(text)
	writer.Close()
	sample.Close()

	r := httptest.NewRequest("POST", "/api/v1/product/upload", body)

	r.Header.Add("Content-Type", writer.FormDataContentType())
	r.AddCookie(&http.Cookie{Name: "session_id", Value: _tmpDB.NewSession("+79990009900")})
	w := httptest.NewRecorder()

	UploadPhotoHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}
