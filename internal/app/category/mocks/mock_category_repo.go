// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category (interfaces: CategoryRepository)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockCategoryRepository is a mock of CategoryRepository interface.
type MockCategoryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCategoryRepositoryMockRecorder
}

// MockCategoryRepositoryMockRecorder is the mock recorder for MockCategoryRepository.
type MockCategoryRepositoryMockRecorder struct {
	mock *MockCategoryRepository
}

// NewMockCategoryRepository creates a new mock instance.
func NewMockCategoryRepository(ctrl *gomock.Controller) *MockCategoryRepository {
	mock := &MockCategoryRepository{ctrl: ctrl}
	mock.recorder = &MockCategoryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCategoryRepository) EXPECT() *MockCategoryRepositoryMockRecorder {
	return m.recorder
}

// SelectCategories mocks base method.
func (m *MockCategoryRepository) SelectCategories() ([]*models.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectCategories")
	ret0, _ := ret[0].([]*models.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectCategories indicates an expected call of SelectCategories.
func (mr *MockCategoryRepositoryMockRecorder) SelectCategories() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectCategories", reflect.TypeOf((*MockCategoryRepository)(nil).SelectCategories))
}
