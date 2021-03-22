package delivery

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session"
	"github.com/gorilla/mux"
	"net/http"
)

type SessionHandler struct {
	sessUcase session.SessionUsecase
}

func NewSessionHandler(sessUcase session.SessionUsecase) *SessionHandler {
	return &SessionHandler{
			sessUcase: sessUcase,
		}
}

func (sh *SessionHandler) Configure(r *mux.Router) {
	r.HandleFunc("/login", sh.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/logout", sh.LogoutHandler).Methods(http.MethodPost)
}

func (sh *SessionHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

}

func (sh *SessionHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {

}