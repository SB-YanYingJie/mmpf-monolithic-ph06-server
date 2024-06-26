package internal

import (
	"context"
	"fmt"
	"testing"

	gi "github.com/machinemapplatform/grpc-interface/golang"
	"github.com/machinemapplatform/library/model"
	im "github.com/machinemapplatform/mmpf-monolithic/internal/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	data1  = []byte("data1")
	data2  = []byte("data2")
	tstamp = &timestamppb.Timestamp{Seconds: 1234, Nanos: 5678}
)

func Test_ValidateSlamRequest(t *testing.T) {
	ctx := context.Background()

	validator := NewGrpcValidator()
	t.Run("ng: no image", func(t *testing.T) {
		input := &gi.SlamRequest{}
		converted, err := validator.ValidateSlamRequest(ctx, input)
		assert.Empty(t, converted)
		assert.EqualError(t, err, "no image")
	})
	t.Run("ng: camera positioin type duplicated", func(t *testing.T) {
		ngDupilicatedTest(t, ctx, validator, gi.LensPlacement_CENTER)
		ngDupilicatedTest(t, ctx, validator, gi.LensPlacement_LEFT)
		ngDupilicatedTest(t, ctx, validator, gi.LensPlacement_RIGHT)
		ngDupilicatedTest(t, ctx, validator, gi.LensPlacement_MERGED)
	})
	t.Run("ng: mismatch image type and camera positioin type", func(t *testing.T) {
		monoNoCenter := &gi.SlamRequest{
			NumberOfLenses: gi.NumberOfLenses_MONO,
			Images: []*gi.Image{
				{LensPlacement: gi.LensPlacement_LEFT, Byte: data1},
			},
		}
		stereoThreeImages := &gi.SlamRequest{
			NumberOfLenses: gi.NumberOfLenses_STEREO,
			Images: []*gi.Image{
				{LensPlacement: gi.LensPlacement_LEFT, Byte: []byte("left")},
				{LensPlacement: gi.LensPlacement_RIGHT, Byte: []byte("right")},
				{LensPlacement: gi.LensPlacement_MERGED, Byte: []byte("both")},
			},
		}
		stereoNoImages := &gi.SlamRequest{
			NumberOfLenses: gi.NumberOfLenses_STEREO,
			Images:         []*gi.Image{},
		}
		stereoBothAndLeftImages := &gi.SlamRequest{
			NumberOfLenses: gi.NumberOfLenses_STEREO,
			Images: []*gi.Image{
				{LensPlacement: gi.LensPlacement_LEFT, Byte: []byte("left")},
				{LensPlacement: gi.LensPlacement_MERGED, Byte: []byte("both")},
			},
		}
		stereoBothAndRightImages := &gi.SlamRequest{
			NumberOfLenses: gi.NumberOfLenses_STEREO,
			Images: []*gi.Image{
				{LensPlacement: gi.LensPlacement_RIGHT, Byte: []byte("right")},
				{LensPlacement: gi.LensPlacement_MERGED, Byte: []byte("both")},
			},
		}
		stereoOnlyLeftImages := &gi.SlamRequest{
			NumberOfLenses: gi.NumberOfLenses_STEREO,
			Images: []*gi.Image{
				{LensPlacement: gi.LensPlacement_LEFT, Byte: []byte("left")},
			},
		}
		stereoOnlyRightImages := &gi.SlamRequest{
			NumberOfLenses: gi.NumberOfLenses_STEREO,
			Images: []*gi.Image{
				{LensPlacement: gi.LensPlacement_RIGHT, Byte: []byte("right")},
			},
		}
		ngMismatchImageAndLensPlacement(t, ctx, validator, monoNoCenter, "no center image was detected, even though image type MONO")
		ngMismatchImageAndLensPlacement(t, ctx, validator, stereoThreeImages, "all of left, right and both image were detected, even though image type STEREO")
		ngMismatchImageAndLensPlacement(t, ctx, validator, stereoNoImages, "any of left, right and both image was detected, even though image type STEREO")
		ngMismatchImageAndLensPlacement(t, ctx, validator, stereoBothAndLeftImages, "left and both image were detected, even though image type STEREO")
		ngMismatchImageAndLensPlacement(t, ctx, validator, stereoBothAndRightImages, "right and both image were detected, even though image type STEREO")
		ngMismatchImageAndLensPlacement(t, ctx, validator, stereoOnlyLeftImages, "no right image was detected, even though image type STEREO")
		ngMismatchImageAndLensPlacement(t, ctx, validator, stereoOnlyRightImages, "no left image was detected, even though image type STEREO")
	})
	t.Run("ng: required param has not detectd", func(t *testing.T) {
		noImageType := &gi.SlamRequest{
			Images: []*gi.Image{
				{LensPlacement: gi.LensPlacement_CENTER, Byte: data1},
			},
			RequestTime: tstamp,
		}
		noImages := &gi.SlamRequest{
			NumberOfLenses: gi.NumberOfLenses_MONO,
			RequestTime:    tstamp,
		}
		noT := &gi.SlamRequest{
			Images: []*gi.Image{
				{LensPlacement: gi.LensPlacement_CENTER, Byte: data1},
			},
			NumberOfLenses: gi.NumberOfLenses_MONO,
		}
		noImagesAndT := &gi.SlamRequest{
			Images: []*gi.Image{
				{LensPlacement: gi.LensPlacement_CENTER, Byte: data1},
			},
		}
		ngRequiredParamNotDetected(t, ctx, validator, noImageType)
		ngRequiredParamNotDetected(t, ctx, validator, noImages)
		ngRequiredParamNotDetected(t, ctx, validator, noT)
		ngRequiredParamNotDetected(t, ctx, validator, noImagesAndT)
	})
	t.Run("ok: ok", func(t *testing.T) {
		monoInput := &gi.SlamRequest{
			NumberOfLenses: gi.NumberOfLenses_MONO,
			Images: []*gi.Image{
				{LensPlacement: gi.LensPlacement_CENTER, Byte: data1},
			},
			RequestTime: tstamp,
		}
		monoExpected := &im.SlamRequest{
			NumberOfLenses: model.MONO,
			Images: map[model.LensPlacement][]byte{
				model.CENTER: data1,
			},
			RequestTime: tstamp.AsTime().UnixNano(),
		}
		stereoLRInput := &gi.SlamRequest{
			NumberOfLenses: gi.NumberOfLenses_STEREO,
			Images: []*gi.Image{
				{LensPlacement: gi.LensPlacement_LEFT, Byte: data1},
				{LensPlacement: gi.LensPlacement_RIGHT, Byte: data2},
			},
			RequestTime: tstamp,
		}
		stereoLRExpected := &im.SlamRequest{
			NumberOfLenses: model.STEREO,
			Images: map[model.LensPlacement][]byte{
				model.LEFT:  data1,
				model.RIGHT: data2,
			},
			RequestTime: tstamp.AsTime().UnixNano(),
		}
		stereoBothInput := &gi.SlamRequest{
			NumberOfLenses: gi.NumberOfLenses_STEREO,
			Images: []*gi.Image{
				{LensPlacement: gi.LensPlacement_MERGED, Byte: data1},
			},
			RequestTime: tstamp,
		}
		stereoBothExpected := &im.SlamRequest{
			NumberOfLenses: model.STEREO,
			Images: map[model.LensPlacement][]byte{
				model.MERGED: data1,
			},
			RequestTime: tstamp.AsTime().UnixNano(),
		}
		ok(t, ctx, validator, monoInput, monoExpected)
		ok(t, ctx, validator, stereoLRInput, stereoLRExpected)
		ok(t, ctx, validator, stereoBothInput, stereoBothExpected)
	})
}

func ngRequiredParamNotDetected(t *testing.T, ctx context.Context, validator *GrpcValidator, input *gi.SlamRequest) {
	actual, err := validator.ValidateSlamRequest(ctx, input)
	assert.Equal(t, (*im.SlamRequest)(nil), actual)
	assert.Error(t, err)
}

func ok(t *testing.T, ctx context.Context, validator *GrpcValidator, input *gi.SlamRequest, expected *im.SlamRequest) {
	actual, err := validator.ValidateSlamRequest(ctx, input)
	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func ngMismatchImageAndLensPlacement(t *testing.T, ctx context.Context, validator *GrpcValidator, input *gi.SlamRequest, errStr string) {
	converted, err := validator.ValidateSlamRequest(ctx, input)
	assert.Empty(t, converted)
	assert.EqualError(t, err, errStr)
}

func ngDupilicatedTest(t *testing.T, ctx context.Context, validator *GrpcValidator, cpt gi.LensPlacement) {
	input := &gi.SlamRequest{
		Images: []*gi.Image{
			{LensPlacement: cpt, Byte: data1},
			{LensPlacement: cpt, Byte: data2},
		},
	}
	converted, err := validator.ValidateSlamRequest(ctx, input)
	assert.Empty(t, converted)
	assert.EqualError(t, err, fmt.Sprintf("duplicated Cameraposition type %s", cpt))
}
