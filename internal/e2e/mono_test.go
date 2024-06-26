//go:build e2e_test
// +build e2e_test

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	gi "github.com/machinemapplatform/grpc-interface/golang"
	"github.com/machinemapplatform/library/asynctest"
	"github.com/machinemapplatform/library/middleware"
	"github.com/machinemapplatform/library/model"
	"github.com/machinemapplatform/library/redis"
	"github.com/machinemapplatform/mmpf-monolithic/cmd/config"
	im "github.com/machinemapplatform/mmpf-monolithic/internal/model"
	"github.com/stretchr/testify/assert"
	"gocv.io/x/gocv"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_Slam_Mono_Success(t *testing.T) {
	var originalRawRequestTime int64 = 1627429514371303100
	ctx := context.Background()
	hostName, _ := os.Hostname()

	r := redis.NewRedis(
		config.RedisAddress,
		config.RedisMaxIdle,
		config.RedisIdleTimeoutSeconds,
		config.RedisPubsubDb,
	)
	defer r.Close()
	deviceClient := NewDeviceClient(address)
	defer deviceClient.Close()
	outChs := []string{"pose_testMMID"}

	pathFormat := fmt.Sprintf("2021_07_27/23/testMMID/raw_image/%d.%%s.png", originalRawRequestTime)
	expectedResponseCPath := fmt.Sprintf(pathFormat, "c")

	t.Run("ok: request slam mono", asynctest.AsyncTest(t, ctx, r, outChs,
		func(t *testing.T) {
			originalRequestId := "testSlamMono"
			originalImagePath := "/app/internal/test_assets/1267_0.png"

			actual, err := deviceClient.C.Slam(ctx, &gi.SlamRequest{
				Metadata:       map[string]string{model.MD_KEY_REQUEST_ID: originalRequestId},
				RequestTime:    timestamppb.New(time.Unix(0, originalRawRequestTime)),
				NumberOfLenses: gi.NumberOfLenses_MONO,
				Images:         []*gi.Image{{LensPlacement: gi.LensPlacement_CENTER, Byte: OpenFileAsBytes(originalImagePath)}},
			})
			assert.NoError(t, err)

			// == responseのチェック ==
			// Metadata
			assert.Equal(t, originalRequestId, actual.Metadata[model.MD_KEY_REQUEST_ID])
			assert.Equal(t, hostName, actual.Metadata[model.MD_KEY_FROM])
			assert.Equal(t, config.MMID, actual.Metadata[model.MD_KEY_MMID])
			assert.Equal(t, strconv.FormatInt(timestamppb.New(time.Unix(0, originalRawRequestTime)).AsTime().UnixNano(), 10), actual.Metadata[model.MD_KEY_REQUEST_TIME])
			assert.Equal(t, config.ImageExt, actual.Metadata[model.MD_KEY_RAW_IMAGE_EXT])
			assert.Equal(t, strconv.Itoa(config.RedisPubsubDb), actual.Metadata[model.MD_KEY_REDIS_IMAGE_DB])
			assert.NotEmpty(t, actual.Metadata[model.MD_KEY_RAW_IMAGE])

			originalImageMat := gocv.IMRead(originalImagePath, gocv.IMReadAnyColor)
			defer originalImageMat.Close()
			originalWidth := originalImageMat.Cols()
			originalHeight := originalImageMat.Rows()
			assert.Equal(t, strconv.Itoa(originalWidth), actual.Metadata[model.MD_KEY_FRAME_SIZE_WIDTH])
			assert.Equal(t, strconv.Itoa(originalHeight), actual.Metadata[model.MD_KEY_FRAME_SIZE_HEIGHT])

			assert.Equal(t, fmt.Sprintf("%.2f", config.Fps), actual.Metadata[model.MD_KEY_TARGET_FPS])

			_, err = strconv.ParseFloat(actual.Metadata[model.MD_KEY_ELAPSED_TIME], 64)
			assert.NoError(t, err)

			// Result
			assert.Equal(t, timestamppb.New(time.Unix(0, originalRawRequestTime)).Seconds, actual.GetPose().RequestTime.Seconds)
			assert.Equal(t, im.ToGrpcSlamState[model.TRACKING_GOOD], actual.GetPose().SlamState)
			assert.Nil(t, actual.GetError())

			// == redisに格納された画像のチェック ==
			actualImage, err := r.BGet(ctx, expectedResponseCPath)
			if err != nil {
				t.Errorf("fail to redis.BGet: %+v", err)
			}
			assert.Equal(t, OpenFileAsBytes(originalImagePath), actualImage)
		},
		func(t *testing.T, channel string, data []byte) {
			originalRequestId := "testSlamMono"
			originalImagePath := "/app/internal/test_assets/1267_0.png"

			// == redisにpublishされたメッセージのチェック ==
			actualData := middleware.DataWithMetadata{}
			json.Unmarshal(data, &actualData)

			// Metadata
			assert.Equal(t, originalRequestId, actualData.Metadata[model.MD_KEY_REQUEST_ID])
			assert.Equal(t, hostName, actualData.Metadata[model.MD_KEY_FROM])
			assert.Equal(t, config.MMID, actualData.Metadata[model.MD_KEY_MMID])
			assert.Equal(t, strconv.FormatInt(timestamppb.New(time.Unix(0, originalRawRequestTime)).AsTime().UnixNano(), 10), actualData.Metadata[model.MD_KEY_REQUEST_TIME])
			assert.Equal(t, config.ImageExt, actualData.Metadata[model.MD_KEY_RAW_IMAGE_EXT])
			assert.Equal(t, strconv.Itoa(config.RedisPubsubDb), actualData.Metadata[model.MD_KEY_REDIS_IMAGE_DB])
			assert.NotEmpty(t, actualData.Metadata[model.MD_KEY_RAW_IMAGE])

			originalImageMat := gocv.IMRead(originalImagePath, gocv.IMReadAnyColor)
			defer originalImageMat.Close()
			originalWidth := originalImageMat.Cols()
			originalHeight := originalImageMat.Rows()
			assert.Equal(t, strconv.Itoa(originalWidth), actualData.Metadata[model.MD_KEY_FRAME_SIZE_WIDTH])
			assert.Equal(t, strconv.Itoa(originalHeight), actualData.Metadata[model.MD_KEY_FRAME_SIZE_HEIGHT])

			assert.Equal(t, fmt.Sprintf("%.2f", config.Fps), actualData.Metadata[model.MD_KEY_TARGET_FPS])

			// Pose
			actualPose := model.Pose{}
			json.Unmarshal(actualData.Data, &actualPose)
			assert.Equal(t, model.TRACKING_GOOD, actualPose.SlamState)
			assert.Equal(t, "pose_testMMID", channel)
		},
	))
}

func Test_Slam_Mono_Failed(t *testing.T) {
	var originalRawRequestTime int64 = 1627429514371303100
	ctx := context.Background()
	hostName, _ := os.Hostname()

	r := redis.NewRedis(
		config.RedisAddress,
		config.RedisMaxIdle,
		config.RedisIdleTimeoutSeconds,
		config.RedisPubsubDb,
	)
	defer r.Close()
	deviceClient := NewDeviceClient(address)
	defer deviceClient.Close()

	t.Run("ng: grpc validation error", func(t *testing.T) {
		originalRequestId := "testSlamMono_grpcValidationError"
		originalImagePath := "/app/internal/test_assets/1267_0.png"

		actual, err := deviceClient.C.Slam(ctx, &gi.SlamRequest{
			Metadata:       map[string]string{model.MD_KEY_REQUEST_ID: originalRequestId},
			RequestTime:    timestamppb.New(time.Unix(0, originalRawRequestTime)),
			NumberOfLenses: gi.NumberOfLenses_STEREO, // 正常系のリクエストではNumberOfLenses_MONO
			Images:         []*gi.Image{{LensPlacement: gi.LensPlacement_CENTER, Byte: OpenFileAsBytes(originalImagePath)}},
		})
		assert.NoError(t, err)

		// == responseのチェック ==
		// Metadata
		assert.Equal(t, originalRequestId, actual.Metadata[model.MD_KEY_REQUEST_ID])

		// Result
		assert.Equal(t, im.ToGrpcErrorStatus[model.INVALID_ARGUMENT], actual.GetError().Status)
		assert.Equal(t, hostName, actual.GetError().From)
		assert.Equal(t, "any of left, right and both image was detected, even though image type STEREO", actual.GetError().Message)
		assert.NotEmpty(t, actual.GetError().ErrorTime.Seconds)
		assert.Nil(t, actual.GetPose())
	})
}
