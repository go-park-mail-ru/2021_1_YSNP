package delivery

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
	"github.com/gorilla/mux"
	"net/http"
)

type UserHandler struct {
	userUcase user.UserUsecase
}

func NewUserHandler(userUcase user.UserUsecase) *UserHandler {
	return &UserHandler{
			userUcase: userUcase,
	}
}

func (uh *UserHandler) Configure(r *mux.Router) {
	r.HandleFunc("/signup", uh.SignUpHandler).Methods(http.MethodPost)
	r.HandleFunc("/upload", uh.UploadAvatarHandler).Methods(http.MethodPost)
	r.HandleFunc("/me", uh.GetProfileHandler).Methods(http.MethodGet)
	r.HandleFunc("/settings", uh.ChangeProfileHandler).Methods(http.MethodPost)
	r.HandleFunc("/settings/password", uh.ChangeProfilePasswordHandler).Methods(http.MethodPost)

}

func (uh *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) ChangeProfileHandler(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) ChangeProfilePasswordHandler(w http.ResponseWriter, r *http.Request) {

}