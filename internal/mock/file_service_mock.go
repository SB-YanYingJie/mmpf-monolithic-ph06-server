// Code generated by MockGen. DO NOT EDIT.
// Source: file_service.go

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	gocv "gocv.io/x/gocv"
)

// MockFileServiceInterface is a mock of FileServiceInterface interface.
type MockFileServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockFileServiceInterfaceMockRecorder
}

// MockFileServiceInterfaceMockRecorder is the mock recorder for MockFileServiceInterface.
type MockFileServiceInterfaceMockRecorder struct {
	mock *MockFileServiceInterface
}

// NewMockFileServiceInterface creates a new mock instance.
func NewMockFileServiceInterface(ctrl *gomock.Controller) *MockFileServiceInterface {
	mock := &MockFileServiceInterface{ctrl: ctrl}
	mock.recorder = &MockFileServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileServiceInterface) EXPECT() *MockFileServiceInterfaceMockRecorder {
	return m.recorder
}

// DecodeFile mocks base method.
func (m *MockFileServiceInterface) DecodeFile(ctx context.Context, data []byte) (*gocv.Mat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecodeFile", ctx, data)
	ret0, _ := ret[0].(*gocv.Mat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DecodeFile indicates an expected call of DecodeFile.
func (mr *MockFileServiceInterfaceMockRecorder) DecodeFile(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecodeFile", reflect.TypeOf((*MockFileServiceInterface)(nil).DecodeFile), ctx, data)
}

// DecodeFiles mocks base method.
func (m *MockFileServiceInterface) DecodeFiles(ctx context.Context, lData, rData []byte) (*gocv.Mat, *gocv.Mat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecodeFiles", ctx, lData, rData)
	ret0, _ := ret[0].(*gocv.Mat)
	ret1, _ := ret[1].(*gocv.Mat)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DecodeFiles indicates an expected call of DecodeFiles.
func (mr *MockFileServiceInterfaceMockRecorder) DecodeFiles(ctx, lData, rData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecodeFiles", reflect.TypeOf((*MockFileServiceInterface)(nil).DecodeFiles), ctx, lData, rData)
}

// WriteFile mocks base method.
func (m *MockFileServiceInterface) WriteFile(ctx context.Context, image []byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteFile", ctx, image)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteFile indicates an expected call of WriteFile.
func (mr *MockFileServiceInterfaceMockRecorder) WriteFile(ctx, image interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFile", reflect.TypeOf((*MockFileServiceInterface)(nil).WriteFile), ctx, image)
}

// WriteFiles mocks base method.
func (m *MockFileServiceInterface) WriteFiles(ctx context.Context, lImage, rImage []byte) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteFiles", ctx, lImage, rImage)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// WriteFiles indicates an expected call of WriteFiles.
func (mr *MockFileServiceInterfaceMockRecorder) WriteFiles(ctx, lImage, rImage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFiles", reflect.TypeOf((*MockFileServiceInterface)(nil).WriteFiles), ctx, lImage, rImage)
}