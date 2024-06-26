package internal

import (
	"context"
	"image"

	"gocv.io/x/gocv"
)

type PreprocessingService struct {
	trimX     int
	trimY     int
	trimWidth int
	trimHight int
}

func NewPreprocessingService(trimX int, trimY int, trimWidth int, trimHight int) *PreprocessingService {
	return &PreprocessingService{
		trimX:     trimX,
		trimY:     trimY,
		trimWidth: trimWidth,
		trimHight: trimHight,
	}
}

func (p *PreprocessingService) ConvertMatToBytes(ctx context.Context, mat *gocv.Mat) []byte {
	return mat.ToBytes()
}

func (p *PreprocessingService) SplitMatIntoTwoMats(ctx context.Context, mat *gocv.Mat) (lByte, rByte *gocv.Mat) {
	lMat, rMat := trim(mat, p.trimX, p.trimY, p.trimWidth, p.trimHight)
	return lMat, rMat
}

func trim(mat *gocv.Mat, x int, y int, width int, height int) (*gocv.Mat, *gocv.Mat) {
	var (
		startPointY = y
		endPointY   = y + height
	)
	var croppedLeft = gocv.NewMat()
	var croppedRight = gocv.NewMat()
	lStartPointX := x
	lEndPointX := x + width
	l := mat.Region(image.Rect(lStartPointX, startPointY, lEndPointX, endPointY))
	defer l.Close()
	l.CopyTo(&croppedLeft)

	rStartPointX := lEndPointX
	rEndPointX := rStartPointX + width
	r := mat.Region(image.Rect(rStartPointX, startPointY, rEndPointX, endPointY))
	defer r.Close()
	r.CopyTo(&croppedRight)
	defer mat.Close()

	return &croppedLeft, &croppedRight
}
