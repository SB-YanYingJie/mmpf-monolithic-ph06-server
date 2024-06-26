package internal

import (
	"context"
	"testing"

	"github.com/KudanJP/KdSlamGo/kdslam"
	"github.com/golang/mock/gomock"
	"github.com/machinemapplatform/library/model"

	mock_internal "github.com/machinemapplatform/mmpf-monolithic/internal/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

type setupTest struct {
	ctx       context.Context
	mmId      string
	calibPath string
	vocabPath string
	kdmpPath  string
	targetFps float32
	mapExpand bool
}

func TestStart(t *testing.T) {
	s := NewSetupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	slamStruct := mock_internal.NewMockSlamInterface(ctrl)
	imageDebugger := mock_internal.NewMockImageDebuggerInterface(ctrl)

	t.Run("start_ok", func(t *testing.T) {
		slamStruct.EXPECT().StartSlam("test/calibPath", "test/vocabPath", float32(30)).Return(nil)
		slamStruct.EXPECT().SetAutoExpansion(false).Return()
		slamStruct.EXPECT().LoadMap("test/kdmapPath", 0, 0).Return(nil)
		slamStruct.EXPECT().GetImageSize().Return(2, 4, nil)
		slamStruct.EXPECT().StopSlam().Return()
		slam := NewSlamService(s.mmId, s.calibPath, s.vocabPath, s.kdmpPath, s.targetFps, false, s.mapExpand, slamStruct, imageDebugger)
		width, height, err := slam.Start(s.ctx)
		assert.Equal(t, 2, width)
		assert.Equal(t, 4, height)
		assert.Equal(t, nil, err)
		slam.Close(s.ctx)
	})

	t.Run("ng_start_error", func(t *testing.T) {
		slamStruct.EXPECT().StartSlam("test/calibPath", "test/vocabPath", float32(30)).Return(xerrors.New("startSlam error"))
		slamStruct.EXPECT().StopSlam().Return()
		slam := NewSlamService(s.mmId, s.calibPath, s.vocabPath, s.kdmpPath, s.targetFps, false, s.mapExpand, slamStruct, imageDebugger)
		width, height, err := slam.Start(s.ctx)
		assert.Equal(t, 0, width)
		assert.Equal(t, 0, height)
		assert.EqualError(t, err, "startSlam error")
		slam.Close(s.ctx)
	})

	t.Run("ng_loadMap_error", func(t *testing.T) {
		slamStruct.EXPECT().StartSlam("test/calibPath", "test/vocabPath", float32(30)).Return(nil)
		slamStruct.EXPECT().SetAutoExpansion(false).Return()
		slamStruct.EXPECT().LoadMap("test/kdmapPath", 0, 0).Return(xerrors.New("LoadMap error"))
		slamStruct.EXPECT().StopSlam().Return()
		slam := NewSlamService(s.mmId, s.calibPath, s.vocabPath, s.kdmpPath, s.targetFps, false, s.mapExpand, slamStruct, imageDebugger)
		width, height, err := slam.Start(s.ctx)
		assert.Equal(t, 0, width)
		assert.Equal(t, 0, height)
		assert.EqualError(t, err, "LoadMap error")
		slam.Close(s.ctx)
	})

	t.Run("ng_GetImageSize_error", func(t *testing.T) {
		slamStruct.EXPECT().StartSlam("test/calibPath", "test/vocabPath", float32(30)).Return(nil)
		slamStruct.EXPECT().SetAutoExpansion(false).Return()
		slamStruct.EXPECT().LoadMap("test/kdmapPath", 0, 0).Return(nil)
		slamStruct.EXPECT().GetImageSize().Return(0, 0, xerrors.New("GetImageSize error"))
		slamStruct.EXPECT().StopSlam().Return()
		slam := NewSlamService(s.mmId, s.calibPath, s.vocabPath, s.kdmpPath, s.targetFps, false, s.mapExpand, slamStruct, imageDebugger)
		width, height, err := slam.Start(s.ctx)
		assert.Equal(t, 0, width)
		assert.Equal(t, 0, height)
		assert.EqualError(t, err, "GetImageSize error")
		slam.Close(s.ctx)
	})
}

func TestGetPoseMono(t *testing.T) {
	s := NewSetupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	slamStruct := mock_internal.NewMockSlamInterface(ctrl)
	imageDebugger := mock_internal.NewMockImageDebuggerInterface(ctrl)

	t.Run("getPoseMono_ok", func(t *testing.T) {

		expectedPose := &kdslam.Pose{
			Tx: 1,
			Ty: 2,
			Tz: 3,
			Qx: 4,
			Qy: 5,
			Qz: 6,
			Qw: 7,
		}

		slamStruct.EXPECT().
			ProcessFrame([]byte("testImage"), false).Return(nil, nil)
		slamStruct.EXPECT().
			GetTransformMatrix().Return(expectedPose, nil)
		slamStruct.EXPECT().
			GetSlamState().Return(model.TRACKING_GOOD)
		imageDebugger.EXPECT().Show(gomock.Any()).Return(nil).Times(0)
		slam := NewSlamService(s.mmId, s.calibPath, s.vocabPath, s.kdmpPath, s.targetFps, false, s.mapExpand, slamStruct, imageDebugger)
		pose, state, err := slam.GetPoseMono(s.ctx, []byte("testImage"))
		t.Log(pose)
		assert.Equal(t, pose, expectedPose)
		assert.Equal(t, model.TRACKING_GOOD, state)
		assert.NoError(t, err)
	})

	t.Run("ng_kdslam_ProcessFrame_error", func(t *testing.T) {

		slamStruct.EXPECT().
			ProcessFrame([]byte("testImage"), false).Return(nil, xerrors.New("processFrame error"))
		imageDebugger.EXPECT().Show(gomock.Any()).Return(nil).Times(0)
		slam := NewSlamService(s.mmId, s.calibPath, s.vocabPath, s.kdmpPath, s.targetFps, false, s.mapExpand, slamStruct, imageDebugger)
		pose, state, err := slam.GetPoseMono(s.ctx, []byte("testImage"))
		t.Log(pose)
		assert.Empty(t, pose)
		assert.Equal(t, model.IDLE, state)
		assert.Error(t, err)
	})

	t.Run("ng_kdslam_GetTransformMatrix_error", func(t *testing.T) {

		slamStruct.EXPECT().
			ProcessFrame([]byte("testImage"), false).Return(nil, nil)
		slamStruct.EXPECT().
			GetTransformMatrix().Return(nil, xerrors.New("getTransformMatrix error"))
		imageDebugger.EXPECT().Show(gomock.Any()).Return(nil).Times(0)
		slam := NewSlamService(s.mmId, s.calibPath, s.vocabPath, s.kdmpPath, s.targetFps, false, s.mapExpand, slamStruct, imageDebugger)
		pose, state, err := slam.GetPoseMono(s.ctx, []byte("testImage"))
		t.Log(pose)
		assert.Empty(t, pose)
		assert.Equal(t, model.IDLE, state)
		assert.Error(t, err)
	})

}
func TestGetPoseStereo(t *testing.T) {
	s := NewSetupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	slamStruct := mock_internal.NewMockSlamInterface(ctrl)
	imageDebugger := mock_internal.NewMockImageDebuggerInterface(ctrl)

	t.Run("ok_pose_data_is_not_empty_if_state_is_2", func(t *testing.T) {
		expectedPose := &kdslam.Pose{
			Tx: 1,
			Ty: 2,
			Tz: 3,
			Qx: 4,
			Qy: 5,
			Qz: 6,
			Qw: 7,
		}
		slamStruct.EXPECT().ProcessFrameStereo([]byte("testImage"), []byte("testImage"), false).Return(nil, nil)
		slamStruct.EXPECT().GetSlamState().Return(model.TRACKING_GOOD)
		slamStruct.EXPECT().GetTransformMatrix().Return(expectedPose, nil)
		imageDebugger.EXPECT().Show(gomock.Any()).Return(nil).Times(0)
		slam := NewSlamService(s.mmId, s.calibPath, s.vocabPath, s.kdmpPath, s.targetFps, false, s.mapExpand, slamStruct, imageDebugger)
		pose, state, err := slam.GetPoseStereo(s.ctx, []byte("testImage"), []byte("testImage"))
		assert.Equal(t, expectedPose, pose)
		assert.Equal(t, model.TRACKING_GOOD, state)
		assert.NoError(t, err)
	})
	t.Run("ok_pose_data_is_empty_if_state_is_not_2", func(t *testing.T) {
		expectedPose := new(kdslam.Pose)
		slamStruct.EXPECT().ProcessFrameStereo([]byte("testImage"), []byte("testImage"), false).Return(nil, nil)
		slamStruct.EXPECT().GetSlamState().Return(model.MAP_LOADING_IN_PROGRESS)
		imageDebugger.EXPECT().Show(gomock.Any()).Return(nil).Times(0)
		slam := NewSlamService(s.mmId, s.calibPath, s.vocabPath, s.kdmpPath, s.targetFps, false, s.mapExpand, slamStruct, imageDebugger)
		pose, state, err := slam.GetPoseStereo(s.ctx, []byte("testImage"), []byte("testImage"))
		assert.Equal(t, expectedPose, pose)
		assert.Equal(t, model.MAP_LOADING_IN_PROGRESS, state)
		assert.NoError(t, err)
	})

	t.Run("ng_ProcessFrameStereo_error", func(t *testing.T) {
		slamStruct.EXPECT().ProcessFrameStereo([]byte("testImage"), []byte("testImage"), false).Return(nil, xerrors.New("ProcessFrameStereo error"))
		imageDebugger.EXPECT().Show(gomock.Any()).Return(nil).Times(0)
		slam := NewSlamService(s.mmId, s.calibPath, s.vocabPath, s.kdmpPath, s.targetFps, false, s.mapExpand, slamStruct, imageDebugger)
		pose, state, err := slam.GetPoseStereo(s.ctx, []byte("testImage"), []byte("testImage"))
		assert.Empty(t, pose)
		assert.Equal(t, model.IDLE, state)
		assert.Error(t, err)
	})

	t.Run("ng_GetTransformMatrix_error", func(t *testing.T) {
		slamStruct.EXPECT().ProcessFrameStereo([]byte("testImage"), []byte("testImage"), false).Return(nil, nil)
		slamStruct.EXPECT().GetSlamState().Return(model.TRACKING_GOOD)
		slamStruct.EXPECT().GetTransformMatrix().Return(nil, xerrors.New("GetTransformMatrix error"))
		imageDebugger.EXPECT().Show(gomock.Any()).Return(nil).Times(0)
		slam := NewSlamService(s.mmId, s.calibPath, s.vocabPath, s.kdmpPath, s.targetFps, false, s.mapExpand, slamStruct, imageDebugger)
		pose, state, err := slam.GetPoseStereo(s.ctx, []byte("testImage"), []byte("testImage"))
		assert.Empty(t, pose)
		assert.Equal(t, model.IDLE, state)
		assert.Error(t, err)

	})

}

func NewSetupTest() setupTest {
	s := setupTest{}
	s.ctx = context.Background()
	s.mmId = "test_mmid"
	s.calibPath = "test/calibPath"
	s.vocabPath = "test/vocabPath"
	s.kdmpPath = "test/kdmapPath"
	s.targetFps = 30
	s.mapExpand = false
	return s
}
