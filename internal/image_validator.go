package internal

import (
	"context"
	"fmt"

	"gocv.io/x/gocv"
)

type ImageValidator struct {
	expectedWidth  int
	expectedHeight int
}

func NewImageValidator(expectedWidth int, expectedHeight int) *ImageValidator {
	return &ImageValidator{
		expectedWidth:  expectedWidth,
		expectedHeight: expectedHeight,
	}
}

func (v *ImageValidator) ValidateImageSize(ctx context.Context, mat *gocv.Mat) error {
	width := mat.Cols()
	if width != v.expectedWidth {
		return fmt.Errorf("invalid image width. expected=%d, actual=%d", v.expectedWidth, width)
	}
	height := mat.Rows()
	if height != v.expectedHeight {
		return fmt.Errorf("invalid image height. expected=%d, actual=%d", v.expectedHeight, height)
	}
	return nil
}
