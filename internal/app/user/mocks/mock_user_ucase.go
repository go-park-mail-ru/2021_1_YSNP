// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user (interfaces: UserUsecase)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	errors "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	models "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUserUsecase is a mock of UserUsecase interface.
type MockUserUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUsecaseMockRecorder
}

// MockUserUsecaseMockRecorder is the mock recorder for MockUserUsecase.
type MockUserUsecaseMockRecorder struct {
	mock *MockUserUsecase
}

// NewMockUserUsecase creates a new mock instance.
func NewMockUserUsecase(ctrl *gomock.Controller) *MockUserUsecase {
	mock := &MockUserUsecase{ctrl: ctrl}
	mock.recorder = &MockUserUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUsecase) EXPECT() *MockUserUsecaseMockRecorder {
	return m.recorder
}

// CheckPassword mocks base method.
func (m *MockUserUsecase) CheckPassword(arg0 *models.UserData, arg1 string) *errors.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPassword", arg0, arg1)
	ret0, _ := ret[0].(*errors.Error)
	return ret0
}

// CheckPassword indicates an expected call of CheckPassword.
func (mr *MockUserUsecaseMockRecorder) CheckPassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPassword", reflect.TypeOf((*MockUserUsecase)(nil).CheckPassword), arg0, arg1)
}

// Create mocks base method.
func (m *MockUserUsecase) Create(arg0 *models.UserData) *errors.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(*errors.Error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUserUsecaseMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserUsecase)(nil).Create), arg0)
}

// GetByID mocks base method.
func (m *MockUserUsecase) GetByID(arg0 uint64) (*models.ProfileData, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0)
	ret0, _ := ret[0].(*models.ProfileData)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockUserUsecaseMockRecorder) GetByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUserUsecase)(nil).GetByID), arg0)
}

// GetByTelephone mocks base method.
func (m *MockUserUsecase) GetByTelephone(arg0 string) (*models.UserData, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTelephone", arg0)
	ret0, _ := ret[0].(*models.UserData)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// GetByTelephone indicates an expected call of GetByTelephone.
func (mr *MockUserUsecaseMockRecorder) GetByTelephone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTelephone", reflect.TypeOf((*MockUserUsecase)(nil).GetByTelephone), arg0)
}

// GetSellerByID mocks base method.
func (m *MockUserUsecase) GetSellerByID(arg0 uint64) (*models.SellerData, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSellerByID", arg0)
	ret0, _ := ret[0].(*models.SellerData)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// GetSellerByID indicates an expected call of GetSellerByID.
func (mr *MockUserUsecaseMockRecorder) GetSellerByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSellerByID", reflect.TypeOf((*MockUserUsecase)(nil).GetSellerByID), arg0)
}

// UpdateAvatar mocks base method.
func (m *MockUserUsecase) UpdateAvatar(arg0 uint64, arg1 string) (*models.UserData, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAvatar", arg0, arg1)
	ret0, _ := ret[0].(*models.UserData)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// UpdateAvatar indicates an expected call of UpdateAvatar.
func (mr *MockUserUsecaseMockRecorder) UpdateAvatar(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAvatar", reflect.TypeOf((*MockUserUsecase)(nil).UpdateAvatar), arg0, arg1)
}

// UpdatePassword mocks base method.
func (m *MockUserUsecase) UpdatePassword(arg0 uint64, arg1 string) (*models.UserData, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", arg0, arg1)
	ret0, _ := ret[0].(*models.UserData)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// UpdatePassword indicates an expected call of UpdatePassword.
func (mr *MockUserUsecaseMockRecorder) UpdatePassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockUserUsecase)(nil).UpdatePassword), arg0, arg1)
}

// UpdatePosition mocks base method.
func (m *MockUserUsecase) UpdatePosition(arg0 uint64, arg1 *models.LocationChangeRequest) (*models.UserData, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLocation", arg0, arg1)
	ret0, _ := ret[0].(*models.UserData)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// UpdatePosition indicates an expected call of UpdatePosition.
func (mr *MockUserUsecaseMockRecorder) UpdatePosition(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLocation", reflect.TypeOf((*MockUserUsecase)(nil).UpdatePosition), arg0, arg1)
}

// UpdateProfile mocks base method.
func (m *MockUserUsecase) UpdateProfile(arg0 uint64, arg1 *models.UserData) (*models.UserData, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfile", arg0, arg1)
	ret0, _ := ret[0].(*models.UserData)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// UpdateProfile indicates an expected call of UpdateProfile.
func (mr *MockUserUsecaseMockRecorder) UpdateProfile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockUserUsecase)(nil).UpdateProfile), arg0, arg1)
}
