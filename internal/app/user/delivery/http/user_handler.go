package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type UserHandler struct {
	userUcase user.UserUsecase
	sessUcase session.SessionUsecase
}

func NewUserHandler(userUcase user.UserUsecase, sessUcase session.SessionUsecase) *UserHandler {
	return &UserHandler{
			userUcase: userUcase,
			sessUcase: sessUcase,
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
	defer r.Body.Close()

	signUp := models.SignUpRequest{}
	err := json.NewDecoder(r.Body).Decode(&signUp)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	user := &models.UserData{
		Name:       signUp.Name,
		Surname:    signUp.Surname,
		Sex:        signUp.Sex,
		Email:      signUp.Email,
		Telephone:  signUp.Telephone,
		Password:   signUp.Password1,
		DateBirth:  signUp.DateBirth,
		LinkImages: signUp.LinkImages,
	}

	err = uh.userUcase.Create(user)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	session := models.CreateSession(user.ID)
	err = uh.sessUcase.Create(session)
	if  err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    session.Value,
		Expires:  session.ExpiresAt,
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful login."))
}

func (uh *UserHandler) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		err := errors.Error{Message: "user not authorised or not found"}
		logrus.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	file, handler, err := r.FormFile("file-upload")
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}
	defer file.Close()
	extension := filepath.Ext(handler.Filename)

	r.FormValue("file-upload")

	str, err := os.Getwd()

	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	photoPath := "static/avatar/"
	os.Chdir(photoPath)

	photoID, err := uuid.NewRandom()
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	f, err := os.OpenFile(photoID.String()+extension, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errors.JSONError(err.Error()))
		return
	}
	defer f.Close()

	os.Chdir(str)

	_, err = io.Copy(f, file)
	if err != nil {
		_ = os.Remove(photoID.String()+extension)
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	avatar := models.Url+"/static/avatar/"+photoID.String()+extension

	_, err = uh.userUcase.UpdateAvatar(userID, avatar)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(avatar)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (uh *UserHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		err := errors.Error{Message: "user not authorised or not found"}
		logrus.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	user, error := uh.userUcase.GetByID(userID)
	if error != nil {
		logrus.Error(error)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(error.Error()))
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (uh *UserHandler) ChangeProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		err := errors.Error{Message: "user not authorised or not found"}
		logrus.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	changeData := models.SignUpRequest{}
	err := json.NewDecoder(r.Body).Decode(&changeData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	user := &models.UserData{
		Name:       changeData.Name,
		Surname:    changeData.Surname,
		Sex:        changeData.Sex,
		Email:      changeData.Email,
		Telephone:  changeData.Telephone,
		Password:   changeData.Password1,
		DateBirth:  changeData.DateBirth,
		LinkImages: changeData.LinkImages,
	}

	_, err = uh.userUcase.UpdateProfile(userID, user)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful change."))
}

func (uh *UserHandler) ChangeProfilePasswordHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		err := errors.Error{Message: "user not authorised or not found"}
		logrus.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	passwordData := models.PasswordChangeRequest{}
	err := json.NewDecoder(r.Body).Decode(&passwordData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	_, err =uh.userUcase.UpdatePassword(userID, passwordData.NewPassword)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful change."))
}