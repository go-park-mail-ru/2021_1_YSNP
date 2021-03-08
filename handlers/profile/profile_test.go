package profile

import (
	_tmpDB "2021_1_YSNP/tmp_database"
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

	r := httptest.NewRequest("POST", "/api/v1/logout", nil)
	r.AddCookie(&http.Cookie{Name:"session_id", Value: _tmpDB.NewSession("+79990009900")})

	var expectedJSON = `{"id":0,"name":"Sergey","surname":"Alehin","sex":"male","email":"alehin@mail.ru","telephone":"+79990009900","password":"Qwerty12","dateBirth":"","linkImages":["http://89.208.199.170:8080/static/avatar/test-avatar.jpg"]}`

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

func TestChangeProfileHandler(t *testing.T) {

}
