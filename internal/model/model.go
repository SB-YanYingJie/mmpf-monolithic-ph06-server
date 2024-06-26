package model

import (
	"github.com/machinemapplatform/library/model"
)

// SlamRequest imitates Request
type SlamRequest struct {
	NumberOfLenses model.NumberOfLenses           `validate:"required"`
	Images         map[model.LensPlacement][]byte `validate:"required"`
	RequestTime    int64                          `validate:"required"`
}
