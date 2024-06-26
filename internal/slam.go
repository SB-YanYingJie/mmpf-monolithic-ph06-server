package internal

import (
	"github.com/KudanJP/KdSlamGo/kdslam"
	"github.com/machinemapplatform/library/model"
	d "github.com/machinemapplatform/mmpf-monolithic/internal/domain"
)

func NewSlam(l kdslam.Logger) d.SlamInterface {
	kdslam.SetLogger(l)
	return &Slam{}
}

type Slam struct{}

func (s *Slam) StartSlam(calibPath string, vocabPath string, fps float32) error {
	return kdslam.StartSlam(calibPath, vocabPath, fps)
}

func (s *Slam) SetAutoExpansion(flag bool) {
	kdslam.SetAutoExpansion(flag)
}

func (s *Slam) LoadMap(mapPath string, syncPolicy int, launchPolicy int) error {
	return kdslam.LoadMap(mapPath, syncPolicy, launchPolicy)
}

func (s *Slam) GetImageSize() (int, int, error) {
	return kdslam.GetImageSize()
}

func (s *Slam) StopSlam() {
	kdslam.StopSlam()
}

func (s *Slam) ProcessFrame(image []byte, debug bool) ([]byte, error) {
	return kdslam.ProcessFrame(image, debug)
}

func (s *Slam) GetTransformMatrix() (*kdslam.Pose, error) {
	return kdslam.GetTransformMatrix()
}

func (s *Slam) GetSlamState() model.SlamState {
	return model.SlamState(kdslam.GetSlamState())
}

func (s *Slam) ProcessFrameStereo(left []byte, right []byte, debug bool) ([]byte, error) {
	return kdslam.ProcessFrameStereo(left, right, debug)
}
