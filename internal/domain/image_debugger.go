package domain

import "gocv.io/x/gocv"

type ImageDebuggerInterface interface {
	Show(mats ...*gocv.Mat) error
}
