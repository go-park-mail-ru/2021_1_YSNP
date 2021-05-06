// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/chat (interfaces: ChatUsecase)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	errors "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	gomock "github.com/golang/mock/gomock"
)

// MockChatUsecase is a mock of ChatUsecase interface.
type MockChatUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockChatUsecaseMockRecorder
}

// MockChatUsecaseMockRecorder is the mock recorder for MockChatUsecase.
type MockChatUsecaseMockRecorder struct {
	mock *MockChatUsecase
}

// NewMockChatUsecase creates a new mock instance.
func NewMockChatUsecase(ctrl *gomock.Controller) *MockChatUsecase {
	mock := &MockChatUsecase{ctrl: ctrl}
	mock.recorder = &MockChatUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatUsecase) EXPECT() *MockChatUsecaseMockRecorder {
	return m.recorder
}

// CreateChat mocks base method.
func (m *MockChatUsecase) CreateChat(arg0 *models.ChatCreateReq, arg1 uint64) (*models.ChatResponse, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChat", arg0, arg1)
	ret0, _ := ret[0].(*models.ChatResponse)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// CreateChat indicates an expected call of CreateChat.
func (mr *MockChatUsecaseMockRecorder) CreateChat(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChat", reflect.TypeOf((*MockChatUsecase)(nil).CreateChat), arg0, arg1)
}

// CreateMessage mocks base method.
func (m *MockChatUsecase) CreateMessage(arg0 *models.CreateMessageReq, arg1 uint64) (*models.MessageResp, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMessage", arg0, arg1)
	ret0, _ := ret[0].(*models.MessageResp)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// CreateMessage indicates an expected call of CreateMessage.
func (mr *MockChatUsecaseMockRecorder) CreateMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMessage", reflect.TypeOf((*MockChatUsecase)(nil).CreateMessage), arg0, arg1)
}

// GetChatById mocks base method.
func (m *MockChatUsecase) GetChatById(arg0, arg1 uint64) (*models.ChatResponse, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChatById", arg0, arg1)
	ret0, _ := ret[0].(*models.ChatResponse)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// GetChatById indicates an expected call of GetChatById.
func (mr *MockChatUsecaseMockRecorder) GetChatById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChatById", reflect.TypeOf((*MockChatUsecase)(nil).GetChatById), arg0, arg1)
}

// GetLastNMessages mocks base method.
func (m *MockChatUsecase) GetLastNMessages(arg0 *models.GetLastNMessagesReq) ([]*models.MessageResp, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastNMessages", arg0)
	ret0, _ := ret[0].([]*models.MessageResp)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// GetLastNMessages indicates an expected call of GetLastNMessages.
func (mr *MockChatUsecaseMockRecorder) GetLastNMessages(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastNMessages", reflect.TypeOf((*MockChatUsecase)(nil).GetLastNMessages), arg0)
}

// GetNMessagesBefore mocks base method.
func (m *MockChatUsecase) GetNMessagesBefore(arg0 *models.GetNMessagesBeforeReq) ([]*models.MessageResp, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNMessagesBefore", arg0)
	ret0, _ := ret[0].([]*models.MessageResp)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// GetNMessagesBefore indicates an expected call of GetNMessagesBefore.
func (mr *MockChatUsecaseMockRecorder) GetNMessagesBefore(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNMessagesBefore", reflect.TypeOf((*MockChatUsecase)(nil).GetNMessagesBefore), arg0)
}

// GetUserChats mocks base method.
func (m *MockChatUsecase) GetUserChats(arg0 uint64) ([]*models.ChatResponse, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserChats", arg0)
	ret0, _ := ret[0].([]*models.ChatResponse)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// GetUserChats indicates an expected call of GetUserChats.
func (mr *MockChatUsecaseMockRecorder) GetUserChats(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserChats", reflect.TypeOf((*MockChatUsecase)(nil).GetUserChats), arg0)
}
