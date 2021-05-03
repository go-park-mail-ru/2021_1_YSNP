package websocket

import (
	"database/sql"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/chat/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/websocket"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestChatWSHandler_CreateMessage_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)
	chatHandler := NewChatWSHandler(chatUcase)

	req := &models.CreateMessageReq{
		ChatID:  0,
		Content: "text",
	}

	var byteData = []byte(`
			{
				"chat_id": 0,
				"content": "text"
				}
	`)

	ctx := &websocket.WSContext{
		Request:  &models.WSMessageReq{
			UserID: 0,
			Type:   "CreateMessageReq",
			Data:   models.CustomData{RequestData: byteData},
		},
	}

	chatUcase.EXPECT().CreateMessage(gomock.Eq(req), ctx.Request.UserID).Return(&models.MessageResp{}, nil)

	chatHandler.CreateMessage(ctx)

	assert.Equal(t, http.StatusOK, ctx.Response.Status)
}

func TestChatWSHandler_CreateMessage_UnmarshErr(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)
	chatHandler := NewChatWSHandler(chatUcase)

	var byteData = []byte(`
			{
				"chat_id": 0,
				"content": "text"
				
	`)

	ctx := &websocket.WSContext{
		Request:  &models.WSMessageReq{
			UserID: 0,
			Type:   "CreateMessageReq",
			Data:   models.CustomData{RequestData: byteData},
		},
	}

	chatHandler.CreateMessage(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Response.Status)
}

func TestChatWSHandler_CreateMessage_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)
	chatHandler := NewChatWSHandler(chatUcase)

	req := &models.CreateMessageReq{
		ChatID:  0,
		Content: "text",
	}

	var byteData = []byte(`
			{
				"chat_id": 0,
				"content": "text"
				}
	`)

	ctx := &websocket.WSContext{
		Request:  &models.WSMessageReq{
			UserID: 0,
			Type:   "CreateMessageReq",
			Data:   models.CustomData{RequestData: byteData},
		},
	}

	chatUcase.EXPECT().CreateMessage(gomock.Eq(req), ctx.Request.UserID).Return(nil, errors.UnexpectedInternal(sql.ErrNoRows))

	chatHandler.CreateMessage(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Response.Status)
}

//func TestChatWSHandler_GetLastNMessages(t *testing.T) {
//	t.Parallel()
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	chatUcase := mock.NewMockChatUsecase(ctrl)
//	chatHandler := NewChatWSHandler(chatUcase)
//
//	req := &models.CreateMessageReq{
//		ChatID:  0,
//		Content: "text",
//	}
//
//	var byteData = []byte(`
//			{
//				"chat_id": 0,
//				"content": "text"
//				}
//	`)
//
//	ctx := &websocket.WSContext{
//		Request:  &models.WSMessageReq{
//			UserID: 0,
//			Type:   "GetLastNMessagesReq",
//			Data:   models.CustomData{RequestData: byteData},
//		},
//	}
//}
