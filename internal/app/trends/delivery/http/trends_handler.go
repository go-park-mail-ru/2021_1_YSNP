package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends"
	"github.com/gorilla/mux"
)

type TrendsHandler struct {
	trendsUsecase trends.TrendsUsecase

}

func NewTrendsHandler(trendsUsecase trends.TrendsUsecase) *TrendsHandler {
	return &TrendsHandler{
		trendsUsecase: trendsUsecase,
	}
}

func (th *TrendsHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/stat", th.LogoutHandler).Methods(http.MethodPost, http.MethodOptions)
}


func (th *TrendsHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	//rawText := "Моя чудесная автомобиль"

	ui := &models.UserInterested{}
	err := json.NewDecoder(r.Body).Decode(&ui)
	if err != nil {
		return
	}
	ui.UserID = 7
	th.trendsUsecase.InsertOrUpdate(ui)
}
