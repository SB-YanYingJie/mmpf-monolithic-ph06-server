package domain

import (
	"context"

	"gocv.io/x/gocv"
)

type ImageValidatorInterface interface {
	ValidateImageSize(ctx context.Context, mat *gocv.Mat) error
}
