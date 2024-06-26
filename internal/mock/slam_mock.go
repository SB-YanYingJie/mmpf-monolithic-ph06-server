// Code generated by MockGen. DO NOT EDIT.
// Source: slam.go

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	reflect "reflect"

	kdslam "github.com/KudanJP/KdSlamGo/kdslam"
	gomock "github.com/golang/mock/gomock"
	model "github.com/machinemapplatform/library/model"
)

// MockSlamInterface is a mock of SlamInterface interface.
type MockSlamInterface struct {
	ctrl     *gomock.Controller
	recorder *MockSlamInterfaceMockRecorder
}

// MockSlamInterfaceMockRecorder is the mock recorder for MockSlamInterface.
type MockSlamInterfaceMockRecorder struct {
	mock *MockSlamInterface
}

// NewMockSlamInterface creates a new mock instance.
func NewMockSlamInterface(ctrl *gomock.Controller) *MockSlamInterface {
	mock := &MockSlamInterface{ctrl: ctrl}
	mock.recorder = &MockSlamInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSlamInterface) EXPECT() *MockSlamInterfaceMockRecorder {
	return m.recorder
}

// GetImageSize mocks base method.
func (m *MockSlamInterface) GetImageSize() (int, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImageSize")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetImageSize indicates an expected call of GetImageSize.
func (mr *MockSlamInterfaceMockRecorder) GetImageSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImageSize", reflect.TypeOf((*MockSlamInterface)(nil).GetImageSize))
}

// GetSlamState mocks base method.
func (m *MockSlamInterface) GetSlamState() model.SlamState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSlamState")
	ret0, _ := ret[0].(model.SlamState)
	return ret0
}

// GetSlamState indicates an expected call of GetSlamState.
func (mr *MockSlamInterfaceMockRecorder) GetSlamState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSlamState", reflect.TypeOf((*MockSlamInterface)(nil).GetSlamState))
}

// GetTransformMatrix mocks base method.
func (m *MockSlamInterface) GetTransformMatrix() (*kdslam.Pose, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransformMatrix")
	ret0, _ := ret[0].(*kdslam.Pose)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransformMatrix indicates an expected call of GetTransformMatrix.
func (mr *MockSlamInterfaceMockRecorder) GetTransformMatrix() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransformMatrix", reflect.TypeOf((*MockSlamInterface)(nil).GetTransformMatrix))
}

// LoadMap mocks base method.
func (m *MockSlamInterface) LoadMap(mapPath string, syncPolicy, launchPolicy int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadMap", mapPath, syncPolicy, launchPolicy)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadMap indicates an expected call of LoadMap.
func (mr *MockSlamInterfaceMockRecorder) LoadMap(mapPath, syncPolicy, launchPolicy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadMap", reflect.TypeOf((*MockSlamInterface)(nil).LoadMap), mapPath, syncPolicy, launchPolicy)
}

// ProcessFrame mocks base method.
func (m *MockSlamInterface) ProcessFrame(image []byte, debug bool) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessFrame", image, debug)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProcessFrame indicates an expected call of ProcessFrame.
func (mr *MockSlamInterfaceMockRecorder) ProcessFrame(image, debug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessFrame", reflect.TypeOf((*MockSlamInterface)(nil).ProcessFrame), image, debug)
}

// ProcessFrameStereo mocks base method.
func (m *MockSlamInterface) ProcessFrameStereo(left, right []byte, debug bool) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessFrameStereo", left, right, debug)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProcessFrameStereo indicates an expected call of ProcessFrameStereo.
func (mr *MockSlamInterfaceMockRecorder) ProcessFrameStereo(left, right, debug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessFrameStereo", reflect.TypeOf((*MockSlamInterface)(nil).ProcessFrameStereo), left, right, debug)
}

// SetAutoExpansion mocks base method.
func (m *MockSlamInterface) SetAutoExpansion(flag bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetAutoExpansion", flag)
}

// SetAutoExpansion indicates an expected call of SetAutoExpansion.
func (mr *MockSlamInterfaceMockRecorder) SetAutoExpansion(flag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAutoExpansion", reflect.TypeOf((*MockSlamInterface)(nil).SetAutoExpansion), flag)
}

// StartSlam mocks base method.
func (m *MockSlamInterface) StartSlam(calibPath, vocabPath string, fps float32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartSlam", calibPath, vocabPath, fps)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartSlam indicates an expected call of StartSlam.
func (mr *MockSlamInterfaceMockRecorder) StartSlam(calibPath, vocabPath, fps interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartSlam", reflect.TypeOf((*MockSlamInterface)(nil).StartSlam), calibPath, vocabPath, fps)
}

// StopSlam mocks base method.
func (m *MockSlamInterface) StopSlam() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StopSlam")
}

// StopSlam indicates an expected call of StopSlam.
func (mr *MockSlamInterfaceMockRecorder) StopSlam() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopSlam", reflect.TypeOf((*MockSlamInterface)(nil).StopSlam))
}