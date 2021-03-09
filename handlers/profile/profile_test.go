package profile

import (
	_tmpDB "2021_1_YSNP/tmp_database"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetProfileHandler_GetProfileHandlerNotAuth(t *testing.T) {
	_tmpDB.InitDB()

	var expectedJSON = `{"message":"User not authorised or not found"}`

	r := httptest.NewRequest("GET", "/api/v1/me", nil)
	w := httptest.NewRecorder()

	GetProfileHandler(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Error("status is not 401")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestGetProfileHandler_GetProfileHandlerSuccess(t *testing.T) {
	_tmpDB.InitDB()

	r := httptest.NewRequest("GET", "/api/v1/me", nil)
	r.AddCookie(&http.Cookie{Name:"session_id", Value: _tmpDB.NewSession("+79990009900")})

	var expectedJSON = `{"id":0,"name":"Sergey","surname":"Alehin","sex":"male","email":"alehin@mail.ru","telephone":"+79990009900","dateBirth":"1991-11-11","linkImages":["http://localhost:8080/static/avatar/test-avatar.jpg"]}`

	w := httptest.NewRecorder()

	GetProfileHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestChangeProfileHandler_ChangeProfileHandlerNotAuth(t *testing.T) {
	_tmpDB.InitDB()

	var expectedJSON = `{"message":"User not authorised or not found"}`

	r := httptest.NewRequest("POST", "/api/v1/setting", nil)
	w := httptest.NewRecorder()

	ChangeProfileHandler(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Error("status is not 401")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}
func TestChangeProfileHandler_ChangeProfileHandlerWrongRequest(t *testing.T) {
	_tmpDB.InitDB()

	expectedJSON := `{"message":"invalid character '}' looking for beginning of object key string"}`

	var byteData = bytes.NewReader([]byte(`{
			"telephone" : "+79990009900",
		}`))

	r := httptest.NewRequest("POST", "/api/v1/setting", byteData)
	r.AddCookie(&http.Cookie{Name:"session_id", Value: _tmpDB.NewSession("+79990009900")})
	w := httptest.NewRecorder()

	ChangeProfileHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not 400")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestChangeProfilePasswordHandler_ChangeProfilePasswordHandlerWrongPass(t *testing.T) {
	var expectedJSON = `{"message":"Old password didn't match."}`

	var byteData = bytes.NewReader([]byte(`{
			"oldPassword" : "Qwerty",
			"newPassword" : "Qwerty12345"
		}`))

	r := httptest.NewRequest("POST", "/api/v1/setting", byteData)
	r.AddCookie(&http.Cookie{Name:"session_id", Value: _tmpDB.NewSession("+79990009900")})
	w := httptest.NewRecorder()

	ChangeProfilePasswordHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("status is not 400")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestChangeProfilePasswordHandler_ChangeProfilePasswordHandlerNoAuth(t *testing.T) {
	var expectedJSON = `{"message":"User not authorised or not found"}`

	r := httptest.NewRequest("POST", "/api/v1/setting", nil)
	w := httptest.NewRecorder()

	ChangeProfilePasswordHandler(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Error("status is not 401")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}

func TestChangeProfilePasswordHandler_ChangeProfilePasswordHandlerSuccess(t *testing.T) {
	var expectedJSON = `{"message":"Successful change."}`

	var byteData = bytes.NewReader([]byte(`{
			"oldPassword" : "Qwerty12",
			"newPassword" : "Qwerty12345"
		}`))

	r := httptest.NewRequest("POST", "/api/v1/setting", byteData)
	r.AddCookie(&http.Cookie{Name:"session_id", Value: _tmpDB.NewSession("+79990009900")})
	w := httptest.NewRecorder()

	ChangeProfilePasswordHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	bytes, _ := ioutil.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}