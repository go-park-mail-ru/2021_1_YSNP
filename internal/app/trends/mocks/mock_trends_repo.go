// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends (interfaces: TrendsRepository)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockTrendsRepository is a mock of TrendsRepository interface.
type MockTrendsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTrendsRepositoryMockRecorder
}

// MockTrendsRepositoryMockRecorder is the mock recorder for MockTrendsRepository.
type MockTrendsRepositoryMockRecorder struct {
	mock *MockTrendsRepository
}

// NewMockTrendsRepository creates a new mock instance.
func NewMockTrendsRepository(ctrl *gomock.Controller) *MockTrendsRepository {
	mock := &MockTrendsRepository{ctrl: ctrl}
	mock.recorder = &MockTrendsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTrendsRepository) EXPECT() *MockTrendsRepositoryMockRecorder {
	return m.recorder
}

// CreateTrendsProducts mocks base method.
func (m *MockTrendsRepository) CreateTrendsProducts(arg0 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTrendsProducts", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTrendsProducts indicates an expected call of CreateTrendsProducts.
func (mr *MockTrendsRepositoryMockRecorder) CreateTrendsProducts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTrendsProducts", reflect.TypeOf((*MockTrendsRepository)(nil).CreateTrendsProducts), arg0)
}

// GetRecommendationProducts mocks base method.
func (m *MockTrendsRepository) GetRecommendationProducts(arg0, arg1 uint64) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecommendationProducts", arg0, arg1)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecommendationProducts indicates an expected call of GetRecommendationProducts.
func (mr *MockTrendsRepositoryMockRecorder) GetRecommendationProducts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecommendationProducts", reflect.TypeOf((*MockTrendsRepository)(nil).GetRecommendationProducts), arg0, arg1)
}

// GetTrendsProducts mocks base method.
func (m *MockTrendsRepository) GetTrendsProducts(arg0 uint64) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrendsProducts", arg0)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrendsProducts indicates an expected call of GetTrendsProducts.
func (mr *MockTrendsRepositoryMockRecorder) GetTrendsProducts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrendsProducts", reflect.TypeOf((*MockTrendsRepository)(nil).GetTrendsProducts), arg0)
}

// InsertOrUpdate mocks base method.
func (m *MockTrendsRepository) InsertOrUpdate(arg0 *models.Trends) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOrUpdate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertOrUpdate indicates an expected call of InsertOrUpdate.
func (mr *MockTrendsRepositoryMockRecorder) InsertOrUpdate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOrUpdate", reflect.TypeOf((*MockTrendsRepository)(nil).InsertOrUpdate), arg0)
}
