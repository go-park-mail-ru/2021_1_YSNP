package signup

import (
	_tmpDB "2021_1_YSNP/tmp_database"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

	expectedJSON := `{"id":0,"name":"Максим","surname":"Торжков","sex":"мужской","email":"a@a.ru","telephone":"+79169230768","password":"Qwerty12","dateBirth":"2021-03-08","linkImages":["http://89.208.199.170:8080/static/avatar/b3c098f5-94d8-4bb9-8e56-bc626e60aab7.jpg"]}`

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

	expectedJSON := `{"message":"User with this phone number exists."}`

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

	//var byteData = bytes.NewReader([]byte(`{"linkImages":"http://89.208.199.170:8080/static/avatar/b3c098f5-94d8-4bb9-8e56-bc626e60aab7.jpg"}`))

	r := httptest.NewRequest("POST", "/api/v1/upload", nil)
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

//func TestUploadAvatarHandler_UploadAvatarHandlerNoFile(t *testing.T) {
//	expectedJSON := `{"message":"request Cofntent-Type isn't multipart/form-data"}`
//
//	//var byteData = bytes.NewReader([]byte(`{"linkImages":"http://89.208.199.170:8080/static/avatar/b3c098f5-94d8-4bb9-8e56-bc626e60aab7.jpg"}`))
//
//	r := httptest.NewRequest("POST", "/upload", nil)
//	r.Header.Add("Content-Type", "multipart/form-data;  boundary=WebKitFormBoundaryWVMPREA66wsYKNBL")
//	w := httptest.NewRecorder()
//
//	UploadAvatarHandler(w, r)
//
//	if w.Code != http.StatusBadRequest {
//		t.Error("Status is not 400")
//	}
//
//	bytes, _ := ioutil.ReadAll(w.Body)
//	if strings.Trim(string(bytes), "\n") != expectedJSON {
//		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
//	}
//}
