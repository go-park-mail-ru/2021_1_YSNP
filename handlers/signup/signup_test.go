package signup

import (
	_tmpDB "2021_1_YSNP/tmp_database"
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestSignUpHandler_SignUpHandlerWrongRequest(t *testing.T) {
	_tmpDB.InitDB()

	expectedJSON := `{"message":"invalid character '}' looking for beginning of object key string"}`

	var byteData = bytes.NewReader([]byte(`{
			"telephone" : "+79990009900",
		}`))

	r := httptest.NewRequest("POST", "/api/v1/signup", byteData)
	w := httptest.NewRecorder()

	SignUpHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestSignUpHandler_SignUpHandlerSucces(t *testing.T) {
	_tmpDB.InitDB()

	expectedJSON := `{"message":"Successful registration."}`

	var byteData = bytes.NewReader([]byte(`{"id":0,
			"name":"Максим",
			"surname":"Торжков",
			"sex":"мужской",
			"email":"a@a.ru",
			"telephone":"+79169230768",
			"password":"Qwerty12",
			"dateBirth":"2021-03-08",
			"linkImages":["http://89.208.199.170:8080/static/avatar/b3c098f5-94d8-4bb9-8e56-bc626e60aab7.jpg"]}`))

	r := httptest.NewRequest("POST", "/api/v1/signup", byteData)
	w := httptest.NewRecorder()

	SignUpHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not 200")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestSignUpHandler_SignUpHandlerUserExists(t *testing.T) {
	_tmpDB.InitDB()

	expectedJSON := `{"message":"user with this phone number exists"}`

	var byteData = bytes.NewReader([]byte(`{"id":0,
			"name":"Максим",
			"surname":"Торжков",
			"sex":"мужской",
			"email":"a@a.ru",
			"telephone":"+79990009900",
			"password":"Qwerty12",
			"dateBirth":"2021-03-08",
			"linkImages":["http://89.208.199.170:8080/static/avatar/b3c098f5-94d8-4bb9-8e56-bc626e60aab7.jpg"]}`))

	r := httptest.NewRequest("POST", "/api/v1/signup", byteData)
	w := httptest.NewRecorder()

	SignUpHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestUploadAvatarHandler_UploadAvatarHandlerWrongContentType(t *testing.T) {
	expectedJSON := `{"message":"request Content-Type isn't multipart/form-data"}`

	r := httptest.NewRequest("POST", "/api/v1/upload", nil)
	r.AddCookie(&http.Cookie{Name: "session_id", Value: _tmpDB.NewSession("+79990009900")})
	w := httptest.NewRecorder()

	UploadAvatarHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestUploadAvatarHandler_UploadAvatarHandlerSucces(t *testing.T) {
	path := "../../static/avatar/test-avatar.jpg"
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file-upload", path)
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

	r := httptest.NewRequest("POST", "/api/v1/upload", body)

	r.Header.Add("Content-Type", writer.FormDataContentType())
	r.AddCookie(&http.Cookie{Name: "session_id", Value: _tmpDB.NewSession("+79990009900")})
	w := httptest.NewRecorder()

	UploadAvatarHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
	}
}

func TestUploadAvatarHandler_UploadAvatarHandlerNoFile(t *testing.T) {
	expectedJSON := `{"message":"http: no such file"}`

	path := "../../static/avatar/test-avatar.jpg"
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", path)
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

	r := httptest.NewRequest("POST", "/api/v1/upload", body)

	r.Header.Add("Content-Type", writer.FormDataContentType())
	r.AddCookie(&http.Cookie{Name: "session_id", Value: _tmpDB.NewSession("+79990009900")})
	w := httptest.NewRecorder()

	UploadAvatarHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestUploadAvatarHandler_UploadAvatarHandlerNotAuth(t *testing.T) {
	path := "../../static/avatar/test-avatar.jpg"
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file-upload", path)
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

	r := httptest.NewRequest("POST", "/api/v1/upload", body)

	r.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	UploadAvatarHandler(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Error("Status is not 401")
	}
}
