package internal

import (
	"context"
	"fmt"

	"github.com/KudanJP/KdSlamGo/kdslam"
	"github.com/machinemapplatform/library/logging"
	"github.com/machinemapplatform/library/model"
	d "github.com/machinemapplatform/mmpf-monolithic/internal/domain"
	"go.uber.org/zap"
	"gocv.io/x/gocv"
)

// NewSlamService is constructor of SlamService.
func NewSlamService(mmID string, calibPath string, vocabPath string, kdmpPath string, targetFps float32, debug bool, mapExpantionFlag bool, slam d.SlamInterface, imageDebugger d.ImageDebuggerInterface) *SlamService {
	return &SlamService{
		mmID:          mmID,
		calib:         calibPath,
		vocab:         vocabPath,
		kdmp:          kdmpPath,
		fps:           targetFps,
		mapExpand:     mapExpantionFlag,
		slam:          slam,
		debug:         debug,
		imageDebugger: imageDebugger,
	}
}

type SlamService struct {
	mmID          string
	calib         string
	vocab         string
	kdmp          string
	fps           float32
	w             int
	h             int
	mapExpand     bool
	slam          d.SlamInterface
	debug         bool
	imageDebugger d.ImageDebuggerInterface
}

func (s *SlamService) Start(ctx context.Context) (imageWidth int, imageHeight int, err error) {

	logger := logging.GetLogger(ctx)

	logger.Debug("SlamService is starting...", zap.String("calib", s.calib),
		zap.String("vocab", s.vocab), zap.Float32("fps", s.fps))

	err = s.slam.StartSlam(s.calib, s.vocab, s.fps)
	if err != nil {
		logger.Error("could not start slam", zap.Error(err))
		return 0, 0, err
	}
	s.slam.SetAutoExpansion(s.mapExpand)

	logger.Debug("SlamService loading map...", zap.String("filename", s.kdmp))
	err = s.slam.LoadMap(s.kdmp, 0, 0)
	if err != nil {
		logger.Error("could not load map", zap.Error(err))
		return 0, 0, err
	}

	s.w, s.h, err = s.slam.GetImageSize()
	if err != nil {
		logger.Error("could not get image size", zap.Error(err))
		return 0, 0, err
	}
	logger.Debug("calibration settings", zap.Int("width", s.w), zap.Int("height", s.h))

	imageWidth = s.w
	imageHeight = s.h
	return imageWidth, imageHeight, nil
}

func (s *SlamService) Close(ctx context.Context) {
	logger := logging.GetLogger(ctx)

	logger.Debug("SlamService is closing...")
	s.slam.StopSlam()
}

func (s *SlamService) GetPoseMono(ctx context.Context, image []byte) (*kdslam.Pose, model.SlamState, error) {
	logger := logging.GetLogger(ctx)

	debug, err := s.slam.ProcessFrame(image, s.debug)
	if err != nil {
		return nil, model.IDLE, fmt.Errorf("kdslam.ProcessFrame failed: %+w", err)
	}
	if s.debug {
		debugMat, _ := gocv.NewMatFromBytes(s.h, s.w, gocv.MatTypeCV8UC4, debug)
		fmt.Println("debug", len(debug))
		fmt.Println("debugMat", debugMat.Size())
		err = s.imageDebugger.Show(&debugMat)
		if err != nil {
			logger.Warn("failed to show debug image", zap.Error(err))
		}
	}
	pose, err := s.slam.GetTransformMatrix()
	if err != nil {
		return nil, model.IDLE, fmt.Errorf("kdslam.GetTransformMatrix failed: %+w", err)
	}
	state := s.slam.GetSlamState()
	return pose, model.SlamState(state), nil
}

func (s *SlamService) GetPoseStereo(ctx context.Context, left []byte, right []byte) (*kdslam.Pose, model.SlamState, error) {
	logger := logging.GetLogger(ctx)

	logger.Debug("image length", zap.Int("Left length", len(left)), zap.Int("Right length", len(right)))
	debug, err := s.slam.ProcessFrameStereo(left, right, s.debug)

	if err != nil {
		return nil, model.IDLE, fmt.Errorf("kdslam.ProcessFrameStereo failed: %+w", err)
	}

	if s.debug {
		debugMat, _ := gocv.NewMatFromBytes(s.h, s.w, gocv.MatTypeCV8UC4, debug)
		logger.Debug("debug image", zap.Int("debug length", len(debug)), zap.Ints("debugMat length", debugMat.Size()))
		err = s.imageDebugger.Show(&debugMat)
		if err != nil {
			logger.Warn("failed to show debug image", zap.Error(err))
		}
	}

	state := s.slam.GetSlamState()
	pose := new(kdslam.Pose)
	if model.SlamState(state) == model.TRACKING_GOOD {
		pose, err = s.slam.GetTransformMatrix()
		if err != nil {
			return nil, model.IDLE, fmt.Errorf("kdslam.GetTransformMatrix failed: %+w", err)
		}
	}

	return pose, model.SlamState(state), nil
}
