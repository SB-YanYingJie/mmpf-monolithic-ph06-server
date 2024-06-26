package domain

import (
	"context"

	gi "github.com/machinemapplatform/grpc-interface/golang"
	im "github.com/machinemapplatform/mmpf-monolithic/internal/model"
)

type GrpcValidatorInterface interface {
	ValidateSlamRequest(ctx context.Context, rawRequest *gi.SlamRequest) (*im.SlamRequest, error)
}
