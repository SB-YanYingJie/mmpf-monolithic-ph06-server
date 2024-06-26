package domain

import (
	"context"

	"github.com/KudanJP/KdSlamGo/kdslam"
	"github.com/machinemapplatform/library/model"
)

type SlamServiceInterface interface {
	Start(ctx context.Context) (imageWidth int, imageHeight int, err error)
	Close(ctx context.Context)
	GetPoseMono(ctx context.Context, image []byte) (*kdslam.Pose, model.SlamState, error)
	GetPoseStereo(ctx context.Context, left []byte, right []byte) (*kdslam.Pose, model.SlamState, error)
}
