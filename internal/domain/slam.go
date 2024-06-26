package domain

import (
	"github.com/KudanJP/KdSlamGo/kdslam"
	"github.com/machinemapplatform/library/model"
)

type SlamInterface interface {
	StartSlam(calibPath string, vocabPath string, fps float32) error
	SetAutoExpansion(flag bool)
	LoadMap(mapPath string, syncPolicy int, launchPolicy int) error
	GetImageSize() (int, int, error)
	StopSlam()
	ProcessFrame(image []byte, debug bool) ([]byte, error)
	GetTransformMatrix() (*kdslam.Pose, error)
	GetSlamState() model.SlamState
	ProcessFrameStereo(left []byte, right []byte, debug bool) ([]byte, error)
}
