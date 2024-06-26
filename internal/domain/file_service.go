package domain

import (
	"context"

	"gocv.io/x/gocv"
)

type FileServiceInterface interface {
	DecodeFile(ctx context.Context, data []byte) (mat *gocv.Mat, err error)
	DecodeFiles(ctx context.Context, lData, rData []byte) (lImage, rImage *gocv.Mat, err error)
	WriteFile(ctx context.Context, image []byte) (string, error)
}
