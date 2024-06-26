package domain

import (
	"context"

	"gocv.io/x/gocv"
)

type PreprocessingServiceInterface interface {
	SplitMatIntoTwoMats(ctx context.Context, mat *gocv.Mat) (lMat, rMat *gocv.Mat)
	ConvertMatToBytes(ctx context.Context, mat *gocv.Mat) []byte
}
