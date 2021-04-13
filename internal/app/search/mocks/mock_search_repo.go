// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search (interfaces: SearchRepository)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockSearchRepository is a mock of SearchRepository interface.
type MockSearchRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSearchRepositoryMockRecorder
}

// MockSearchRepositoryMockRecorder is the mock recorder for MockSearchRepository.
type MockSearchRepositoryMockRecorder struct {
	mock *MockSearchRepository
}

// NewMockSearchRepository creates a new mock instance.
func NewMockSearchRepository(ctrl *gomock.Controller) *MockSearchRepository {
	mock := &MockSearchRepository{ctrl: ctrl}
	mock.recorder = &MockSearchRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSearchRepository) EXPECT() *MockSearchRepositoryMockRecorder {
	return m.recorder
}

// SelectByFilter mocks base method.
func (m *MockSearchRepository) SelectByFilter(arg0 *uint64, arg1 *models.Search) ([]*models.ProductListData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectByFilter", arg0, arg1)
	ret0, _ := ret[0].([]*models.ProductListData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectByFilter indicates an expected call of SelectByFilter.
func (mr *MockSearchRepositoryMockRecorder) SelectByFilter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectByFilter", reflect.TypeOf((*MockSearchRepository)(nil).SelectByFilter), arg0, arg1)
}
