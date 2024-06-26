// Code generated by MockGen. DO NOT EDIT.
// Source: /app/internal/domain/preprocessing_service.go

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	gocv "gocv.io/x/gocv"
)

// MockPreprocessingServiceInterface is a mock of PreprocessingServiceInterface interface.
type MockPreprocessingServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockPreprocessingServiceInterfaceMockRecorder
}

// MockPreprocessingServiceInterfaceMockRecorder is the mock recorder for MockPreprocessingServiceInterface.
type MockPreprocessingServiceInterfaceMockRecorder struct {
	mock *MockPreprocessingServiceInterface
}

// NewMockPreprocessingServiceInterface creates a new mock instance.
func NewMockPreprocessingServiceInterface(ctrl *gomock.Controller) *MockPreprocessingServiceInterface {
	mock := &MockPreprocessingServiceInterface{ctrl: ctrl}
	mock.recorder = &MockPreprocessingServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPreprocessingServiceInterface) EXPECT() *MockPreprocessingServiceInterfaceMockRecorder {
	return m.recorder
}

// ConvertMatToBytes mocks base method.
func (m *MockPreprocessingServiceInterface) ConvertMatToBytes(ctx context.Context, mat *gocv.Mat) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConvertMatToBytes", ctx, mat)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// ConvertMatToBytes indicates an expected call of ConvertMatToBytes.
func (mr *MockPreprocessingServiceInterfaceMockRecorder) ConvertMatToBytes(ctx, mat interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConvertMatToBytes", reflect.TypeOf((*MockPreprocessingServiceInterface)(nil).ConvertMatToBytes), ctx, mat)
}

// SplitMatIntoTwoMats mocks base method.
func (m *MockPreprocessingServiceInterface) SplitMatIntoTwoMats(ctx context.Context, mat *gocv.Mat) (*gocv.Mat, *gocv.Mat) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SplitMatIntoTwoMats", ctx, mat)
	ret0, _ := ret[0].(*gocv.Mat)
	ret1, _ := ret[1].(*gocv.Mat)
	return ret0, ret1
}

// SplitMatIntoTwoMats indicates an expected call of SplitMatIntoTwoMats.
func (mr *MockPreprocessingServiceInterfaceMockRecorder) SplitMatIntoTwoMats(ctx, mat interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SplitMatIntoTwoMats", reflect.TypeOf((*MockPreprocessingServiceInterface)(nil).SplitMatIntoTwoMats), ctx, mat)
}
