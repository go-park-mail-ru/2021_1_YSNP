// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search (interfaces: SearchUsecase)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	errors "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	gomock "github.com/golang/mock/gomock"
)

// MockSearchUsecase is a mock of SearchUsecase interface.
type MockSearchUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockSearchUsecaseMockRecorder
}

// MockSearchUsecaseMockRecorder is the mock recorder for MockSearchUsecase.
type MockSearchUsecaseMockRecorder struct {
	mock *MockSearchUsecase
}

// NewMockSearchUsecase creates a new mock instance.
func NewMockSearchUsecase(ctrl *gomock.Controller) *MockSearchUsecase {
	mock := &MockSearchUsecase{ctrl: ctrl}
	mock.recorder = &MockSearchUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSearchUsecase) EXPECT() *MockSearchUsecaseMockRecorder {
	return m.recorder
}

// SelectByFilter mocks base method.
func (m *MockSearchUsecase) SelectByFilter(arg0 *uint64, arg1 *models.Search) ([]*models.ProductListData, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectByFilter", arg0, arg1)
	ret0, _ := ret[0].([]*models.ProductListData)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// SelectByFilter indicates an expected call of SelectByFilter.
func (mr *MockSearchUsecaseMockRecorder) SelectByFilter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectByFilter", reflect.TypeOf((*MockSearchUsecase)(nil).SelectByFilter), arg0, arg1)
}
