package login

import (
	_tmpDB "2021_1_YSNP/tmp_database"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginHandler_LoginHandlerSuccess(t *testing.T) {
	_tmpDB.InitDB()

	var byteData = bytes.NewReader([]byte(`{
			"telephone" : "+79990009900",
			"password" : "Qwerty12"
		}`))

	r := httptest.NewRequest("POST", "/login", byteData)
	w := httptest.NewRecorder()

	LoginHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
	}

}

func TestLoginHandler_LoginHandlerWrongRequest(t *testing.T) {
	_tmpDB.InitDB()

	expectedJSON := `{"message":"invalid character '}' looking for beginning of object key string"}`

	var byteData = bytes.NewReader([]byte(`{
			"telephone" : "+79990009900",
			"password" : "Qwerty12",
		}`))

	r := httptest.NewRequest("POST", "/login", byteData)
	w := httptest.NewRecorder()

	LoginHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}

}

func TestLoginHandler_LoginHandlerNoUser(t *testing.T) {
	_tmpDB.InitDB()

	expectedJSON := `{"message":"No user with this number"}`

	var byteData = bytes.NewReader([]byte(`{
			"telephone" : "+7",
			"password" : "Qwerty12"
		}`))

	r := httptest.NewRequest("POST", "/login", byteData)
	w := httptest.NewRecorder()

	LoginHandler(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("Status is not 404")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestLogoutHandler(t *testing.T) {

}
