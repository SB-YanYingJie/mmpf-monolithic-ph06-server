package internal

import (
	"fmt"

	"gocv.io/x/gocv"
)

type ImageDebuggerInterface interface {
	Show(mats ...*gocv.Mat) error
}

type ImageDebugger struct {
	windows    []*gocv.Window
	windowsLen int
}

func NewImageDebugger(debug bool, windowNames ...string) (*ImageDebugger, error) {
	var debugWindows = []*gocv.Window{}
	if debug {
		for _, name := range windowNames {
			if name == "" {
				return nil, fmt.Errorf("window name must be nonempty")
			}
			debugWindows = append(debugWindows, gocv.NewWindow(name))
		}
	}
	return &ImageDebugger{
		windows:    debugWindows,
		windowsLen: len(debugWindows),
	}, nil
}

func (d *ImageDebugger) Show(mats ...*gocv.Mat) error {
	if d.windowsLen != len(mats) {
		return fmt.Errorf("mats length must be equivalent to windows length")
	}

	for i, mat := range mats {
		d.windows[i].IMShow(*mat)
		d.windows[i].WaitKey(1)
	}
	return nil
}
