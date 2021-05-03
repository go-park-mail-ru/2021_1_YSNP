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

func TestChatWSHandler_GetLastNMessages_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)
	chatHandler := NewChatWSHandler(chatUcase)

	req := &models.GetLastNMessagesReq{
		UserID: 0,
		ChatID:  0,
		Count: 1,
	}

	var byteData = []byte(`
			{
				"chat_id": 0,
				"count": 1
				}
	`)

	ctx := &websocket.WSContext{
		Request:  &models.WSMessageReq{
			UserID: 0,
			Type:   "GetLastNMessagesReq",
			Data:   models.CustomData{RequestData: byteData},
		},
	}

	chatUcase.EXPECT().GetLastNMessages(req).Return([]*models.MessageResp{}, nil)

	chatHandler.GetLastNMessages(ctx)

	assert.Equal(t, http.StatusOK, ctx.Response.Status)
}

func TestChatWSHandler_GetLastNMessages_UnmarshErr(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)
	chatHandler := NewChatWSHandler(chatUcase)

	var byteData = []byte(`
			{
				"chat_id": 0,
				"count": 1
				
	`)

	ctx := &websocket.WSContext{
		Request:  &models.WSMessageReq{
			UserID: 0,
			Type:   "GetLastNMessagesReq",
			Data:   models.CustomData{RequestData: byteData},
		},
	}

	chatHandler.GetLastNMessages(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Response.Status)
}

func TestChatWSHandler_GetLastNMessages_Err(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)
	chatHandler := NewChatWSHandler(chatUcase)

	req := &models.GetLastNMessagesReq{
		UserID: 0,
		ChatID:  0,
		Count: 1,
	}

	var byteData = []byte(`
			{
				"chat_id": 0,
				"count": 1
				}
	`)

	ctx := &websocket.WSContext{
		Request:  &models.WSMessageReq{
			UserID: 0,
			Type:   "GetLastNMessagesReq",
			Data:   models.CustomData{RequestData: byteData},
		},
	}

	chatUcase.EXPECT().GetLastNMessages(req).Return(nil, errors.UnexpectedInternal(sql.ErrNoRows))

	chatHandler.GetLastNMessages(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Response.Status)
}

func TestChatWSHandler_GetNMessagesBefore_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)
	chatHandler := NewChatWSHandler(chatUcase)

	req := &models.GetNMessagesBeforeReq{
		ChatID:  0,
		Count: 1,
		LastMessageID: 0,
	}

	var byteData = []byte(`
			{
				"chat_id": 0,
				"count": 1,
				"message_id": 0
				}
	`)

	ctx := &websocket.WSContext{
		Request:  &models.WSMessageReq{
			UserID: 0,
			Type:   "GetNMessagesBeforeReq",
			Data:   models.CustomData{RequestData: byteData},
		},
	}

	chatUcase.EXPECT().GetNMessagesBefore(req).Return([]*models.MessageResp{}, nil)

	chatHandler.GetNMessagesBefore(ctx)

	assert.Equal(t, http.StatusOK, ctx.Response.Status)
}

func TestChatWSHandler_GetNMessagesBefore_UnmarshErr(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)
	chatHandler := NewChatWSHandler(chatUcase)


	var byteData = []byte(`
			{
				"chat_id": 0,
				"count": 1,
				"message_id": 0
				
	`)

	ctx := &websocket.WSContext{
		Request:  &models.WSMessageReq{
			UserID: 0,
			Type:   "GetNMessagesBeforeReq",
			Data:   models.CustomData{RequestData: byteData},
		},
	}

	chatHandler.GetNMessagesBefore(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Response.Status)
}

func TestChatWSHandler_GetNMessagesBefore_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)
	chatHandler := NewChatWSHandler(chatUcase)

	req := &models.GetNMessagesBeforeReq{
		ChatID:  0,
		Count: 1,
		LastMessageID: 0,
	}

	var byteData = []byte(`
			{
				"chat_id": 0,
				"count": 1,
				"message_id": 0
				}
	`)

	ctx := &websocket.WSContext{
		Request:  &models.WSMessageReq{
			UserID: 0,
			Type:   "GetNMessagesBeforeReq",
			Data:   models.CustomData{RequestData: byteData},
		},
	}

	chatUcase.EXPECT().GetNMessagesBefore(req).Return(nil, errors.UnexpectedInternal(sql.ErrNoRows))

	chatHandler.GetNMessagesBefore(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Response.Status)
}