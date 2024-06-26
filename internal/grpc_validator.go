package internal

import (
	"context"
	"fmt"

	"github.com/go-playground/validator"
	gi "github.com/machinemapplatform/grpc-interface/golang"
	"github.com/machinemapplatform/library/model"
	im "github.com/machinemapplatform/mmpf-monolithic/internal/model"
)

type GrpcValidator struct {
	validate *validator.Validate
}

func NewGrpcValidator() *GrpcValidator {
	v := validator.New()
	return &GrpcValidator{validate: v}
}

// ValidateSlamRequest validate grpc request. if no error, returns converted request
func (v GrpcValidator) ValidateSlamRequest(ctx context.Context, rawRequest *gi.SlamRequest) (*im.SlamRequest, error) {

	// 同一のカメラポジションが複数含まれないこと
	images, err := isCamPosTypeUnique(rawRequest.Images)
	if err != nil {
		return nil, err
	}

	// イメージタイプとカメラポジションの関係性が正しいこと
	imageType := im.ToModelImageType[rawRequest.NumberOfLenses]
	err = validateCamPosTypeAndImageTypeRelation(imageType, images)
	if err != nil {
		return nil, err
	}

	request := &im.SlamRequest{
		NumberOfLenses: imageType,
		Images:         images,
		RequestTime:    rawRequest.RequestTime.AsTime().UnixNano(),
	}

	err = v.validate.Struct(request)
	if err != nil {
		return nil, fmt.Errorf("invalid argument: %+w", err)
	}

	return request, nil
}

// model.MONOの場合、CENTER の画像があれば true、そうでなければ false を返却します。
// model.STEREOの場合、BOTH の 画像のみ、または（LEFT、RIGHT）の画像があれば true、そうでなければ false を返却します。
func validateCamPosTypeAndImageTypeRelation(imageType model.NumberOfLenses, images map[model.LensPlacement][]byte) error {
	switch imageType {
	case model.MONO:
		if images[model.CENTER] == nil {
			return fmt.Errorf("no center image was detected, even though image type MONO")
		}
	case model.STEREO:
		l := images[model.LEFT]
		r := images[model.RIGHT]
		m := images[model.MERGED]

		if m != nil && l != nil && r != nil {
			return fmt.Errorf("all of left, right and both image were detected, even though image type STEREO")
		}
		if m == nil && l == nil && r == nil {
			return fmt.Errorf("any of left, right and both image was detected, even though image type STEREO")
		}

		if m != nil {
			if l != nil {
				return fmt.Errorf("left and both image were detected, even though image type STEREO")
			}
			if r != nil {
				return fmt.Errorf("right and both image were detected, even though image type STEREO")
			}
		} else {
			if l == nil {
				return fmt.Errorf("no left image was detected, even though image type STEREO")
			}
			if r == nil {
				return fmt.Errorf("no right image was detected, even though image type STEREO")
			}
		}
	}
	return nil
}

func isCamPosTypeUnique(images []*gi.Image) (map[model.LensPlacement][]byte, error) {
	if images == nil {
		return nil, fmt.Errorf("no image")
	}
	mapped := map[model.LensPlacement][]byte{}
	counter := map[model.LensPlacement]int{
		model.CENTER: 0,
		model.RIGHT:  0,
		model.LEFT:   0,
		model.MERGED: 0,
	}
	for _, v := range images {
		switch v.LensPlacement {
		case gi.LensPlacement_CENTER:
			if counter[model.CENTER]++; counter[model.CENTER] > 1 {
				return nil, fmt.Errorf("duplicated Cameraposition type %s", gi.LensPlacement_CENTER)
			}
			mapped[model.CENTER] = v.Byte
		case gi.LensPlacement_LEFT:
			if counter[model.LEFT]++; counter[model.LEFT] > 1 {
				return nil, fmt.Errorf("duplicated Cameraposition type %s", gi.LensPlacement_LEFT)
			}
			mapped[model.LEFT] = v.Byte
		case gi.LensPlacement_RIGHT:
			if counter[model.RIGHT]++; counter[model.RIGHT] > 1 {
				return nil, fmt.Errorf("duplicated Cameraposition type %s", gi.LensPlacement_RIGHT)
			}
			mapped[model.RIGHT] = v.Byte
		case gi.LensPlacement_MERGED:
			if counter[model.MERGED]++; counter[model.MERGED] > 1 {
				return nil, fmt.Errorf("duplicated Cameraposition type %s", gi.LensPlacement_MERGED)
			}
			mapped[model.MERGED] = v.Byte
		}
	}
	return mapped, nil
}
