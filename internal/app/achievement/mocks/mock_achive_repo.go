// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2021_1_YSNP/internal/app/achievement (interfaces: AchievementRepository)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockAchievementRepository is a mock of AchievementRepository interface.
type MockAchievementRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAchievementRepositoryMockRecorder
}

// MockAchievementRepositoryMockRecorder is the mock recorder for MockAchievementRepository.
type MockAchievementRepositoryMockRecorder struct {
	mock *MockAchievementRepository
}

// NewMockAchievementRepository creates a new mock instance.
func NewMockAchievementRepository(ctrl *gomock.Controller) *MockAchievementRepository {
	mock := &MockAchievementRepository{ctrl: ctrl}
	mock.recorder = &MockAchievementRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAchievementRepository) EXPECT() *MockAchievementRepositoryMockRecorder {
	return m.recorder
}

// GetUserAchievements mocks base method.
func (m *MockAchievementRepository) GetUserAchievements(arg0 int) ([]*models.Achievement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserAchievements", arg0)
	ret0, _ := ret[0].([]*models.Achievement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserAchievements indicates an expected call of GetUserAchievements.
func (mr *MockAchievementRepositoryMockRecorder) GetUserAchievements(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserAchievements", reflect.TypeOf((*MockAchievementRepository)(nil).GetUserAchievements), arg0)
}
