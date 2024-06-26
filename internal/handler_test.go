package internal

import (
	"context"
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/KudanJP/KdSlamGo/kdslam"
	gomock "github.com/golang/mock/gomock"
	gi "github.com/machinemapplatform/grpc-interface/golang"
	"github.com/machinemapplatform/library/model"
	"github.com/machinemapplatform/library/mytime"
	rm "github.com/machinemapplatform/library/redis/mock"
	"github.com/machinemapplatform/mmpf-monolithic/cmd/config"
	d "github.com/machinemapplatform/mmpf-monolithic/internal/domain"
	mock "github.com/machinemapplatform/mmpf-monolithic/internal/mock"
	im "github.com/machinemapplatform/mmpf-monolithic/internal/model"
	"github.com/machinemapplatform/mmpf-monolithic/pkg"
	"github.com/stretchr/testify/assert"
	"gocv.io/x/gocv"
)

func Test_Slam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	gv := mock.NewMockGrpcValidatorInterface(ctrl)
	fs := mock.NewMockFileServiceInterface(ctrl)
	iv := mock.NewMockImageValidatorInterface(ctrl)
	ps := mock.NewMockPreprocessingServiceInterface(ctrl)
	ss := mock.NewMockSlamServiceInterface(ctrl)
	rm := rm.NewMockRedisInterface(ctrl)

	timeInt := 12345678
	ctx := mytime.WithTime(context.Background(), int64(timeInt))
	logger := pkg.InitLogger(config.LogSettingsFilePath, config.ServiceName)
	ctx = pkg.WithLogger(ctx, logger)
	hostname, _ := os.Hostname()

	t.Run("ok_slam_mono", func(t *testing.T) {
		tmp := config.ImageType
		config.ImageType = "mono"
		defer func() {
			config.ImageType = tmp
		}()
		requestId := "unused_client"
		metadata := map[string]string{
			model.MD_KEY_MMID:       config.MMID,
			model.MD_KEY_REQUEST_ID: requestId,
		}
		rawImage := []byte("unused image")
		rawRequest := &gi.SlamRequest{
			Metadata:       metadata,
			RequestTime:    toTimestamppb(mytime.NowUnixNano(ctx)),
			NumberOfLenses: gi.NumberOfLenses_MONO,
			Images: []*gi.Image{
				{
					LensPlacement: gi.LensPlacement_CENTER,
					Byte:          rawImage,
				},
			},
		}
		validSr := &im.SlamRequest{
			NumberOfLenses: model.MONO,
			Images:         map[model.LensPlacement][]byte{model.CENTER: rawImage},
			RequestTime:    int64(timeInt),
		}
		path := "unused path"
		kdpose := &kdslam.Pose{
			Tx: 2.2345678,
			Ty: 2.3345678,
			Tz: 2.4345678,
			Qx: 2.4345678,
			Qy: 2.5345678,
			Qz: 2.6345678,
			Qw: 2.7345678,
		}
		state := model.SlamState(gi.SlamState_TRACKING_GOOD)
		pose := createPubPose(kdpose, state)
		mat := &gocv.Mat{}

		prImg := []byte("preprocessed image")

		gv.EXPECT().ValidateSlamRequest(gomock.Any(), rawRequest).Return(validSr, nil)
		fs.EXPECT().WriteFile(gomock.Any(), rawImage).Return(path, nil)
		fs.EXPECT().DecodeFile(gomock.Any(), rawImage).Return(mat, nil)
		iv.EXPECT().ValidateImageSize(gomock.Any(), mat).Return(nil)
		ps.EXPECT().ConvertMatToBytes(gomock.Any(), mat).Return(prImg)
		ss.EXPECT().GetPoseMono(gomock.Any(), prImg).Return(kdpose, state, nil)
		rm.EXPECT().Publish(gomock.Any(), config.RedisPubsubPoseChannel, gomock.Any()).Return(0, nil)

		height := "680"
		width := "400"
		dbNum := "2"
		fpsStr := "30"
		h := NewHandler(
			d.GrpcConnectorField{
				GrpcValidator: gv,
			},
			d.PreprocessField{
				FrameSizeWidth:       width,
				FrameSizeHeight:      height,
				ImageDbNumber:        dbNum,
				FileService:          fs,
				ImageValidator:       iv,
				PreprocessingService: ps,
			},
			d.SlamField{
				FpsStr:      fpsStr,
				SlamService: ss,
				Redis:       rm,
			},
		)
		excmetadata := safeMergeMap(nil, metadata, "client_")
		excmetadata[model.MD_KEY_REQUEST_ID] = requestId
		excmetadata[model.MD_KEY_MMID] = config.MMID
		excmetadata[model.MD_KEY_REQUEST_TIME] = strconv.Itoa(timeInt)
		// endtime - starttime = 0になるため、0でアサートする。
		excmetadata[model.MD_KEY_ELAPSED_TIME] = "0.000"
		excmetadata[model.MD_KEY_RAW_IMAGE] = path
		excmetadata[model.MD_KEY_RAW_IMAGE_EXT] = config.ImageExt
		excmetadata[model.MD_KEY_FRAME_SIZE_WIDTH] = width
		excmetadata[model.MD_KEY_FRAME_SIZE_HEIGHT] = height
		excmetadata[model.MD_KEY_REDIS_IMAGE_DB] = dbNum
		excmetadata[model.MD_KEY_TARGET_FPS] = fpsStr
		excmetadata[model.MD_KEY_FROM] = hostname

		expRes := &gi.SlamResponse{
			Metadata: excmetadata,
		}
		expRes.Result = &gi.SlamResponse_Pose{
			Pose: &gi.PoseResult{
				RequestTime: toTimestamppb(int64(timeInt)),
				PosX:        pose.PosX,
				PosY:        pose.PosY,
				PosZ:        pose.PosZ,
				QuatX:       pose.QuatX,
				QuatY:       pose.QuatY,
				QuatZ:       pose.QuatZ,
				QuatW:       pose.QuatW,
				SlamState:   im.ToGrpcSlamState[state],
			},
		}

		actRes, actErr := h.Slam(ctx, rawRequest)
		assert.Nil(t, actErr)
		assert.Equal(t, expRes, actRes)

	})

	t.Run("ok_slam_stereo(separated)", func(t *testing.T) {
		tmp := config.ImageType
		config.ImageType = "stereo_separated"
		defer func() {
			config.ImageType = tmp
		}()

		requestId := "unused_client"
		metadata := map[string]string{
			model.MD_KEY_MMID:       config.MMID,
			model.MD_KEY_REQUEST_ID: requestId,
		}
		rawlImage := []byte("unused limage")
		rawrImage := []byte("unused rimage")
		rawRequest := &gi.SlamRequest{
			Metadata:       metadata,
			RequestTime:    toTimestamppb(mytime.NowUnixNano(ctx)),
			NumberOfLenses: gi.NumberOfLenses_STEREO,
			Images: []*gi.Image{
				{
					LensPlacement: gi.LensPlacement_LEFT,
					Byte:          rawlImage,
				},
				{
					LensPlacement: gi.LensPlacement_RIGHT,
					Byte:          rawrImage,
				},
			},
		}
		validSr := &im.SlamRequest{
			NumberOfLenses: model.STEREO,
			Images: map[model.LensPlacement][]byte{
				model.LEFT:  rawlImage,
				model.RIGHT: rawrImage,
			},
			RequestTime: int64(timeInt),
		}
		path := "unused path"
		kdpose := &kdslam.Pose{
			Tx: 2.2345678,
			Ty: 2.3345678,
			Tz: 2.4345678,
			Qx: 2.4345678,
			Qy: 2.5345678,
			Qz: 2.6345678,
			Qw: 2.7345678,
		}
		state := model.SlamState(gi.SlamState_TRACKING_GOOD)
		pose := createPubPose(kdpose, state)
		mat := &gocv.Mat{}

		prImg := []byte("preprocessed image")

		gv.EXPECT().ValidateSlamRequest(gomock.Any(), rawRequest).Return(validSr, nil)
		fs.EXPECT().WriteFile(gomock.Any(), rawlImage).Return(path, nil)
		fs.EXPECT().DecodeFiles(gomock.Any(), rawlImage, rawrImage).Return(mat, mat, nil)
		iv.EXPECT().ValidateImageSize(gomock.Any(), mat).Return(nil).Times(2)
		ps.EXPECT().ConvertMatToBytes(gomock.Any(), mat).Return(prImg).Times(2)
		ss.EXPECT().GetPoseStereo(gomock.Any(), prImg, prImg).Return(kdpose, state, nil)
		rm.EXPECT().Publish(gomock.Any(), config.RedisPubsubPoseChannel, gomock.Any()).Return(0, nil)

		height := "680"
		width := "400"
		dbNum := "2"
		fpsStr := "30"
		h := NewHandler(
			d.GrpcConnectorField{
				GrpcValidator: gv,
			},
			d.PreprocessField{
				FrameSizeWidth:       width,
				FrameSizeHeight:      height,
				ImageDbNumber:        dbNum,
				FileService:          fs,
				ImageValidator:       iv,
				PreprocessingService: ps,
			},
			d.SlamField{
				FpsStr:      fpsStr,
				SlamService: ss,
				Redis:       rm,
			},
		)
		excmetadata := safeMergeMap(nil, metadata, "client_")
		excmetadata[model.MD_KEY_REQUEST_ID] = requestId
		excmetadata[model.MD_KEY_MMID] = config.MMID
		excmetadata[model.MD_KEY_REQUEST_TIME] = strconv.Itoa(timeInt)
		// endtime - starttime = 0になるため、0でアサートする。
		excmetadata[model.MD_KEY_ELAPSED_TIME] = "0.000"
		excmetadata[model.MD_KEY_RAW_IMAGE] = path
		excmetadata[model.MD_KEY_RAW_IMAGE_EXT] = config.ImageExt
		excmetadata[model.MD_KEY_FRAME_SIZE_WIDTH] = width
		excmetadata[model.MD_KEY_FRAME_SIZE_HEIGHT] = height
		excmetadata[model.MD_KEY_REDIS_IMAGE_DB] = dbNum
		excmetadata[model.MD_KEY_TARGET_FPS] = fpsStr
		excmetadata[model.MD_KEY_FROM] = hostname

		expRes := &gi.SlamResponse{
			Metadata: excmetadata,
		}
		expRes.Result = &gi.SlamResponse_Pose{
			Pose: &gi.PoseResult{
				RequestTime: toTimestamppb(int64(timeInt)),
				PosX:        pose.PosX,
				PosY:        pose.PosY,
				PosZ:        pose.PosZ,
				QuatX:       pose.QuatX,
				QuatY:       pose.QuatY,
				QuatZ:       pose.QuatZ,
				QuatW:       pose.QuatW,
				SlamState:   im.ToGrpcSlamState[state],
			},
		}

		actRes, actErr := h.Slam(ctx, rawRequest)
		assert.Nil(t, actErr)
		assert.Equal(t, expRes, actRes)
	})

	t.Run("ok_slam_stereo(merged)", func(t *testing.T) {
		tmp := config.ImageType
		config.ImageType = "stereo_merged"
		defer func() {
			config.ImageType = tmp
		}()

		requestId := "unused_client"
		metadata := map[string]string{
			model.MD_KEY_MMID:       config.MMID,
			model.MD_KEY_REQUEST_ID: requestId,
		}
		rawImage := []byte("unused image")
		rawRequest := &gi.SlamRequest{
			Metadata:       metadata,
			RequestTime:    toTimestamppb(mytime.NowUnixNano(ctx)),
			NumberOfLenses: gi.NumberOfLenses_STEREO,
			Images: []*gi.Image{
				{
					LensPlacement: gi.LensPlacement_MERGED,
					Byte:          rawImage,
				},
			},
		}
		validSr := &im.SlamRequest{
			NumberOfLenses: model.STEREO,
			Images:         map[model.LensPlacement][]byte{model.MERGED: rawImage},
			RequestTime:    int64(timeInt),
		}
		path := "unused path"
		kdpose := &kdslam.Pose{
			Tx: 2.2345678,
			Ty: 2.3345678,
			Tz: 2.4345678,
			Qx: 2.4345678,
			Qy: 2.5345678,
			Qz: 2.6345678,
			Qw: 2.7345678,
		}
		state := model.SlamState(gi.SlamState_TRACKING_GOOD)
		pose := createPubPose(kdpose, state)
		mat := &gocv.Mat{}

		expMatL := &gocv.Mat{}
		expMatR := &gocv.Mat{}

		prRImg := []byte("preprocessed rimage")
		prLImg := []byte("preprocessed limage")

		gv.EXPECT().ValidateSlamRequest(gomock.Any(), rawRequest).Return(validSr, nil)
		fs.EXPECT().WriteFile(gomock.Any(), rawImage).Return(path, nil)
		fs.EXPECT().DecodeFile(gomock.Any(), rawImage).Return(mat, nil)
		ps.EXPECT().SplitMatIntoTwoMats(gomock.Any(), mat).Return(expMatL, expMatR)
		iv.EXPECT().ValidateImageSize(gomock.Any(), expMatL).Return(nil)
		iv.EXPECT().ValidateImageSize(gomock.Any(), expMatR).Return(nil)
		ps.EXPECT().ConvertMatToBytes(gomock.Any(), expMatL).Return(prLImg)
		ps.EXPECT().ConvertMatToBytes(gomock.Any(), expMatR).Return(prRImg)
		ss.EXPECT().GetPoseStereo(gomock.Any(), prLImg, prRImg).Return(kdpose, state, nil)
		rm.EXPECT().Publish(gomock.Any(), config.RedisPubsubPoseChannel, gomock.Any()).Return(0, nil)

		height := "680"
		width := "400"
		dbNum := "2"
		fpsStr := "30"
		h := NewHandler(
			d.GrpcConnectorField{
				GrpcValidator: gv,
			},
			d.PreprocessField{
				FrameSizeWidth:       width,
				FrameSizeHeight:      height,
				ImageDbNumber:        dbNum,
				FileService:          fs,
				ImageValidator:       iv,
				PreprocessingService: ps,
			},
			d.SlamField{
				FpsStr:      fpsStr,
				SlamService: ss,
				Redis:       rm,
			},
		)
		excmetadata := safeMergeMap(nil, metadata, "client_")
		excmetadata[model.MD_KEY_REQUEST_ID] = requestId
		excmetadata[model.MD_KEY_MMID] = config.MMID
		excmetadata[model.MD_KEY_REQUEST_TIME] = strconv.Itoa(timeInt)
		// endtime - starttime = 0になるため、0でアサートする。
		excmetadata[model.MD_KEY_ELAPSED_TIME] = "0.000"
		excmetadata[model.MD_KEY_RAW_IMAGE] = path
		excmetadata[model.MD_KEY_RAW_IMAGE_EXT] = config.ImageExt
		excmetadata[model.MD_KEY_FRAME_SIZE_WIDTH] = width
		excmetadata[model.MD_KEY_FRAME_SIZE_HEIGHT] = height
		excmetadata[model.MD_KEY_REDIS_IMAGE_DB] = dbNum
		excmetadata[model.MD_KEY_TARGET_FPS] = fpsStr
		excmetadata[model.MD_KEY_FROM] = hostname

		expRes := &gi.SlamResponse{
			Metadata: excmetadata,
		}
		expRes.Result = &gi.SlamResponse_Pose{
			Pose: &gi.PoseResult{
				RequestTime: toTimestamppb(int64(timeInt)),
				PosX:        pose.PosX,
				PosY:        pose.PosY,
				PosZ:        pose.PosZ,
				QuatX:       pose.QuatX,
				QuatY:       pose.QuatY,
				QuatZ:       pose.QuatZ,
				QuatW:       pose.QuatW,
				SlamState:   im.ToGrpcSlamState[state],
			},
		}

		actRes, actErr := h.Slam(ctx, rawRequest)
		assert.Nil(t, actErr)
		assert.Equal(t, expRes, actRes)
	})

	t.Run("ng_slam_request_id_not_found_and_invalid_argument", func(t *testing.T) {
		metadata := map[string]string{
			model.MD_KEY_MMID: config.MMID,
		}
		rawImage := []byte("unused image")
		rawRequest := &gi.SlamRequest{
			Metadata:       metadata,
			RequestTime:    toTimestamppb(mytime.NowUnixNano(ctx)),
			NumberOfLenses: gi.NumberOfLenses_MONO,
			Images: []*gi.Image{
				{
					LensPlacement: gi.LensPlacement_CENTER,
					Byte:          rawImage,
				},
			},
		}
		err := errors.New("invalid argument")
		gv.EXPECT().ValidateSlamRequest(gomock.Any(), rawRequest).Return(nil, err)

		excmetadata := metadata

		expRes := &gi.SlamResponse{
			Metadata: excmetadata,
		}
		expRes.Result = &gi.SlamResponse_Error{
			Error: &gi.ErrorResult{
				Status:    im.ToGrpcErrorStatus[model.INVALID_ARGUMENT],
				Message:   err.Error(),
				From:      hostname,
				ErrorTime: toTimestamppb(int64(timeInt)),
			},
		}

		h := NewHandler(
			d.GrpcConnectorField{
				GrpcValidator: gv,
			},
			d.PreprocessField{},
			d.SlamField{},
		)

		actRes, actErr := h.Slam(ctx, rawRequest)
		assert.Error(t, err, actErr)
		assert.Equal(t, expRes, actRes)
	})

	t.Run("ng_slam_invalid_env_imageType", func(t *testing.T) {
		tmp := config.ImageType
		config.ImageType = "hoge"
		defer func() {
			config.ImageType = tmp
		}()

		requestId := "unused_client"
		metadata := map[string]string{
			model.MD_KEY_MMID:       config.MMID,
			model.MD_KEY_REQUEST_ID: requestId,
		}
		rawImage := []byte("unused image")
		rawRequest := &gi.SlamRequest{
			Metadata:       metadata,
			RequestTime:    toTimestamppb(mytime.NowUnixNano(ctx)),
			NumberOfLenses: gi.NumberOfLenses_MONO,
			Images: []*gi.Image{
				{
					LensPlacement: gi.LensPlacement_CENTER,
					Byte:          rawImage,
				},
			},
		}
		validSr := &im.SlamRequest{
			NumberOfLenses: model.MONO,
			Images:         map[model.LensPlacement][]byte{model.CENTER: rawImage},
			RequestTime:    int64(timeInt),
		}

		gv.EXPECT().ValidateSlamRequest(gomock.Any(), rawRequest).Return(validSr, nil)

		h := NewHandler(
			d.GrpcConnectorField{
				GrpcValidator: gv,
			},
			d.PreprocessField{},
			d.SlamField{},
		)

		assert.Panics(t, func() {
			_, err := h.Slam(ctx, rawRequest)
			if err != nil {
				t.Log(err)
			}
		})
	})

	t.Run("ng_slam_failed_to_get_slam_result", func(t *testing.T) {
		tmp := config.ImageType
		config.ImageType = "mono"
		defer func() {
			config.ImageType = tmp
		}()
		requestId := "unused_client"
		metadata := map[string]string{
			model.MD_KEY_MMID:       config.MMID,
			model.MD_KEY_REQUEST_ID: requestId,
		}
		rawImage := []byte("unused image")
		rawRequest := &gi.SlamRequest{
			Metadata:       metadata,
			RequestTime:    toTimestamppb(mytime.NowUnixNano(ctx)),
			NumberOfLenses: gi.NumberOfLenses_MONO,
			Images: []*gi.Image{
				{
					LensPlacement: gi.LensPlacement_CENTER,
					Byte:          rawImage,
				},
			},
		}
		validSr := &im.SlamRequest{
			NumberOfLenses: model.MONO,
			Images:         map[model.LensPlacement][]byte{model.CENTER: rawImage},
			RequestTime:    int64(timeInt),
		}
		path := "unused path"

		gv.EXPECT().ValidateSlamRequest(gomock.Any(), rawRequest).Return(validSr, nil)
		fs.EXPECT().WriteFile(gomock.Any(), rawImage).Return(path, nil)
		err := errors.New("failed to write")
		fs.EXPECT().DecodeFile(gomock.Any(), rawImage).Return(nil, err)

		h := NewHandler(
			d.GrpcConnectorField{
				GrpcValidator: gv,
			},
			d.PreprocessField{
				FileService: fs,
			},
			d.SlamField{},
		)
		excmetadata := metadata

		expRes := &gi.SlamResponse{
			Metadata: excmetadata,
		}
		expRes.Result = &gi.SlamResponse_Error{
			Error: &gi.ErrorResult{
				Status:    im.ToGrpcErrorStatus[model.INTERNAL],
				Message:   err.Error(),
				From:      hostname,
				ErrorTime: toTimestamppb(int64(timeInt)),
			},
		}

		actRes, actErr := h.Slam(ctx, rawRequest)
		assert.Error(t, err, actErr)
		assert.Equal(t, expRes, actRes)
	})

	t.Run("ng_slam_failed_to_publish_pose", func(t *testing.T) {
		tmp := config.ImageType
		config.ImageType = "mono"
		defer func() {
			config.ImageType = tmp
		}()
		requestId := "unused_client"
		metadata := map[string]string{
			model.MD_KEY_MMID:       config.MMID,
			model.MD_KEY_REQUEST_ID: requestId,
		}
		rawImage := []byte("unused image")
		rawRequest := &gi.SlamRequest{
			Metadata:       metadata,
			RequestTime:    toTimestamppb(mytime.NowUnixNano(ctx)),
			NumberOfLenses: gi.NumberOfLenses_MONO,
			Images: []*gi.Image{
				{
					LensPlacement: gi.LensPlacement_CENTER,
					Byte:          rawImage,
				},
			},
		}
		validSr := &im.SlamRequest{
			NumberOfLenses: model.MONO,
			Images:         map[model.LensPlacement][]byte{model.CENTER: rawImage},
			RequestTime:    int64(timeInt),
		}
		path := "unused path"
		kdpose := &kdslam.Pose{
			Tx: 2.2345678,
			Ty: 2.3345678,
			Tz: 2.4345678,
			Qx: 2.4345678,
			Qy: 2.5345678,
			Qz: 2.6345678,
			Qw: 2.7345678,
		}
		state := model.SlamState(gi.SlamState_TRACKING_GOOD)
		pose := createPubPose(kdpose, state)
		mat := &gocv.Mat{}

		prImg := []byte("preprocessed image")

		gv.EXPECT().ValidateSlamRequest(gomock.Any(), rawRequest).Return(validSr, nil)
		fs.EXPECT().WriteFile(gomock.Any(), rawImage).Return(path, nil)
		fs.EXPECT().DecodeFile(gomock.Any(), rawImage).Return(mat, nil)
		iv.EXPECT().ValidateImageSize(gomock.Any(), mat).Return(nil)
		ps.EXPECT().ConvertMatToBytes(gomock.Any(), mat).Return(prImg)
		ss.EXPECT().GetPoseMono(gomock.Any(), prImg).Return(kdpose, state, nil)
		err := errors.New("failed to publish")
		rm.EXPECT().Publish(gomock.Any(), config.RedisPubsubPoseChannel, gomock.Any()).Return(0, err)

		// sample.iniを想定
		height := "680"
		width := "400"
		dbNum := "2"
		fpsStr := "30"
		h := NewHandler(
			d.GrpcConnectorField{
				GrpcValidator: gv,
			},
			d.PreprocessField{
				FrameSizeWidth:       width,
				FrameSizeHeight:      height,
				ImageDbNumber:        dbNum,
				FileService:          fs,
				ImageValidator:       iv,
				PreprocessingService: ps,
			},
			d.SlamField{
				FpsStr:      fpsStr,
				SlamService: ss,
				Redis:       rm,
			},
		)
		excmetadata := safeMergeMap(nil, metadata, "client_")
		excmetadata[model.MD_KEY_REQUEST_ID] = requestId
		excmetadata[model.MD_KEY_MMID] = config.MMID
		excmetadata[model.MD_KEY_REQUEST_TIME] = strconv.Itoa(timeInt)
		// endtime - starttime = 0になるため、0でアサートする。
		excmetadata[model.MD_KEY_ELAPSED_TIME] = "0.000"
		excmetadata[model.MD_KEY_RAW_IMAGE] = path
		excmetadata[model.MD_KEY_RAW_IMAGE_EXT] = config.ImageExt
		excmetadata[model.MD_KEY_FRAME_SIZE_WIDTH] = width
		excmetadata[model.MD_KEY_FRAME_SIZE_HEIGHT] = height
		excmetadata[model.MD_KEY_REDIS_IMAGE_DB] = dbNum
		excmetadata[model.MD_KEY_TARGET_FPS] = fpsStr
		excmetadata[model.MD_KEY_FROM] = hostname

		expRes := &gi.SlamResponse{
			Metadata: excmetadata,
		}
		expRes.Result = &gi.SlamResponse_Pose{
			Pose: &gi.PoseResult{
				RequestTime: toTimestamppb(int64(timeInt)),
				PosX:        pose.PosX,
				PosY:        pose.PosY,
				PosZ:        pose.PosZ,
				QuatX:       pose.QuatX,
				QuatY:       pose.QuatY,
				QuatZ:       pose.QuatZ,
				QuatW:       pose.QuatW,
				SlamState:   im.ToGrpcSlamState[state],
			},
		}

		actRes, actErr := h.Slam(ctx, rawRequest)
		assert.Nil(t, actErr)
		assert.Equal(t, expRes, actRes)

	})

}

func Test_PreProcessMono(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fs := mock.NewMockFileServiceInterface(ctrl)
	ps := mock.NewMockPreprocessingServiceInterface(ctrl)
	iv := mock.NewMockImageValidatorInterface(ctrl)

	ctx := context.Background()
	logger := pkg.InitLogger(config.LogSettingsFilePath, config.ServiceName)
	ctx = pkg.WithLogger(ctx, logger)

	image := []byte("unused image")

	path := "unused"
	mat := &gocv.Mat{}

	t.Run("ok_preprocess_mono", func(t *testing.T) {
		expByte := []byte("expect image byte")
		fs.EXPECT().WriteFile(ctx, image).Return(path, nil)
		fs.EXPECT().DecodeFile(ctx, image).Return(mat, nil)
		iv.EXPECT().ValidateImageSize(ctx, mat).Return(nil)
		ps.EXPECT().ConvertMatToBytes(ctx, mat).Return(expByte)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{
				PreprocessingService: ps,
				ImageValidator:       iv,
				FileService:          fs,
				Debug:                false,
			},
			d.SlamField{},
		)

		actByte, actStr, err := h.PreProcessMono(ctx, image)
		assert.Equal(t, nil, err)
		assert.Equal(t, expByte, actByte)
		assert.Equal(t, path, actStr)
	})

	t.Run("ng_preprocess_mono_failed_writefile_and_decodefile", func(t *testing.T) {
		wErr := errors.New("failed to write")
		fs.EXPECT().WriteFile(ctx, image).Return("", wErr)
		dErr := errors.New("failed to decode")
		fs.EXPECT().DecodeFile(ctx, image).Return(nil, dErr)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{
				PreprocessingService: ps,
				ImageValidator:       iv,
				FileService:          fs,
				Debug:                true,
			},
			d.SlamField{},
		)

		actByte, actStr, actErr := h.PreProcessMono(ctx, image)
		assert.Equal(t, dErr, actErr)
		assert.Nil(t, actByte)
		assert.Empty(t, actStr)
	})

	t.Run("ng_preprocess_mono_failed_show_and_validate", func(t *testing.T) {
		id := mock.NewMockImageDebuggerInterface(ctrl)

		fs.EXPECT().WriteFile(ctx, image).Return(path, nil)
		fs.EXPECT().DecodeFile(ctx, image).Return(mat, nil)
		sErr := errors.New("failed to show")
		id.EXPECT().Show(mat).Return(sErr)
		uErr := errors.New("unMatch image size")
		iv.EXPECT().ValidateImageSize(ctx, mat).Return(uErr)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{
				PreprocessingService: ps,
				ImageValidator:       iv,
				FileService:          fs,
				Debug:                true,
				ImageDebugger:        id,
			},
			d.SlamField{},
		)

		actByte, actStr, actErr := h.PreProcessMono(ctx, image)
		assert.Equal(t, uErr, actErr)
		assert.Nil(t, actByte)
		assert.Empty(t, actStr)
	})
}

func Test_PreProcessStereoSeparated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fs := mock.NewMockFileServiceInterface(ctrl)
	ps := mock.NewMockPreprocessingServiceInterface(ctrl)
	iv := mock.NewMockImageValidatorInterface(ctrl)

	ctx := context.Background()
	logger := pkg.InitLogger(config.LogSettingsFilePath, config.ServiceName)
	ctx = pkg.WithLogger(ctx, logger)

	rimage := []byte("unused rimage")
	limage := []byte("unused limage")

	lpath := "unused l"
	rmat := &gocv.Mat{}
	lmat := &gocv.Mat{}

	t.Run("ok_preprocess_streo(separated)", func(t *testing.T) {
		explByte := []byte("expect image lbyte")
		exprByte := []byte("expect image tbyte")
		fs.EXPECT().WriteFile(ctx, limage).Return(lpath, nil)
		fs.EXPECT().DecodeFiles(ctx, limage, rimage).Return(lmat, rmat, nil)
		iv.EXPECT().ValidateImageSize(ctx, lmat).Return(nil)
		iv.EXPECT().ValidateImageSize(ctx, rmat).Return(nil)
		ps.EXPECT().ConvertMatToBytes(ctx, lmat).Return(explByte)
		ps.EXPECT().ConvertMatToBytes(ctx, rmat).Return(exprByte)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{
				PreprocessingService: ps,
				ImageValidator:       iv,
				FileService:          fs,
				Debug:                false,
			},
			d.SlamField{},
		)

		actlByte, actrByte, actStr, err := h.PreProcessStereoSeparated(ctx, limage, rimage)
		assert.Equal(t, nil, err)
		assert.Equal(t, explByte, actlByte)
		assert.Equal(t, exprByte, actrByte)
		assert.Equal(t, lpath, actStr)
	})

	t.Run("ng_preprocess_stereo(separated)_failed_writefile_and_decodefile", func(t *testing.T) {
		wErr := errors.New("failed to write")
		fs.EXPECT().WriteFile(ctx, limage).Return(lpath, wErr)
		dErr := errors.New("failed to decode")
		fs.EXPECT().DecodeFiles(ctx, limage, rimage).Return(nil, nil, dErr)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{
				PreprocessingService: ps,
				ImageValidator:       iv,
				FileService:          fs,
				Debug:                true,
			},
			d.SlamField{},
		)

		actlByte, actrByte, actStr, actErr := h.PreProcessStereoSeparated(ctx, limage, rimage)
		assert.Equal(t, dErr, actErr)
		assert.Nil(t, actlByte)
		assert.Nil(t, actrByte)
		assert.Empty(t, actStr)
	})

	t.Run("ng_preprocess_stereo(separated)_failed_show_and_validate_l", func(t *testing.T) {
		id := mock.NewMockImageDebuggerInterface(ctrl)

		fs.EXPECT().WriteFile(ctx, limage).Return(lpath, nil)
		fs.EXPECT().DecodeFiles(ctx, limage, rimage).Return(lmat, rmat, nil)
		sErr := errors.New("failed to show")
		id.EXPECT().Show(lmat).Return(sErr)
		uErr := errors.New("unMatch image size")
		iv.EXPECT().ValidateImageSize(ctx, lmat).Return(uErr)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{
				PreprocessingService: ps,
				ImageValidator:       iv,
				FileService:          fs,
				Debug:                true,
				ImageDebugger:        id,
			},
			d.SlamField{},
		)

		actlByte, actrByte, actStr, actErr := h.PreProcessStereoSeparated(ctx, limage, rimage)
		assert.Equal(t, uErr, actErr)
		assert.Nil(t, actlByte)
		assert.Nil(t, actrByte)
		assert.Empty(t, actStr)
	})

	t.Run("ng_preprocess_stereo(separated)_failed_show_and_validate_r", func(t *testing.T) {
		id := mock.NewMockImageDebuggerInterface(ctrl)

		fs.EXPECT().WriteFile(ctx, limage).Return(lpath, nil)
		fs.EXPECT().DecodeFiles(ctx, limage, rimage).Return(lmat, rmat, nil)
		sErr := errors.New("failed to show")
		id.EXPECT().Show(lmat).Return(sErr)
		uErr := errors.New("unMatch image size")
		iv.EXPECT().ValidateImageSize(ctx, lmat).Return(nil)
		iv.EXPECT().ValidateImageSize(ctx, rmat).Return(uErr)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{
				PreprocessingService: ps,
				ImageValidator:       iv,
				FileService:          fs,
				Debug:                true,
				ImageDebugger:        id,
			},
			d.SlamField{},
		)

		actlByte, actrByte, actStr, actErr := h.PreProcessStereoSeparated(ctx, limage, rimage)
		assert.Equal(t, uErr, actErr)
		assert.Nil(t, actlByte)
		assert.Nil(t, actrByte)
		assert.Empty(t, actStr)
	})
}

func Test_PreProcessStreoMerged(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fs := mock.NewMockFileServiceInterface(ctrl)
	ps := mock.NewMockPreprocessingServiceInterface(ctrl)
	iv := mock.NewMockImageValidatorInterface(ctrl)

	ctx := context.Background()
	logger := pkg.InitLogger(config.LogSettingsFilePath, config.ServiceName)
	ctx = pkg.WithLogger(ctx, logger)

	image := []byte("unused image")

	path := "unused"
	mat := &gocv.Mat{}

	t.Run("ok_preprocess_stereo(merged)", func(t *testing.T) {
		expMatL := &gocv.Mat{}
		expMatR := &gocv.Mat{}

		explByte := []byte("expect limage byte")
		exprByte := []byte("expect rimage byte")

		fs.EXPECT().WriteFile(ctx, image).Return(path, nil)
		fs.EXPECT().DecodeFile(ctx, image).Return(mat, nil)
		ps.EXPECT().SplitMatIntoTwoMats(ctx, mat).Return(expMatL, expMatR)
		iv.EXPECT().ValidateImageSize(ctx, expMatL).Return(nil)
		iv.EXPECT().ValidateImageSize(ctx, expMatR).Return(nil)
		ps.EXPECT().ConvertMatToBytes(ctx, expMatL).Return(explByte)
		ps.EXPECT().ConvertMatToBytes(ctx, expMatR).Return(exprByte)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{
				PreprocessingService: ps,
				ImageValidator:       iv,
				FileService:          fs,
				Debug:                false,
			},
			d.SlamField{},
		)

		actlByte, actrByte, actStr, err := h.PreProcessStereoMerged(ctx, image)
		assert.Equal(t, nil, err)
		assert.Equal(t, explByte, actlByte)
		assert.Equal(t, exprByte, actrByte)
		assert.Equal(t, path, actStr)
	})

	t.Run("ng_preprocess_stereo(merged)_failed_writefile_and_decodefile", func(t *testing.T) {
		wErr := errors.New("failed to write")
		fs.EXPECT().WriteFile(ctx, image).Return("", wErr)
		dErr := errors.New("failed to decode")
		fs.EXPECT().DecodeFile(ctx, image).Return(nil, dErr)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{
				PreprocessingService: ps,
				ImageValidator:       iv,
				FileService:          fs,
				Debug:                true,
			},
			d.SlamField{},
		)

		actlByte, actrByte, actStr, actErr := h.PreProcessStereoMerged(ctx, image)
		assert.Equal(t, dErr, actErr)
		assert.Nil(t, actlByte)
		assert.Nil(t, actrByte)
		assert.Empty(t, actStr)
	})

	t.Run("ng_preprocess_stereo(merged)_failed_show_and_validate", func(t *testing.T) {
		expMatL := &gocv.Mat{}
		expMatR := &gocv.Mat{}
		id := mock.NewMockImageDebuggerInterface(ctrl)
		fs.EXPECT().WriteFile(ctx, image).Return(path, nil)
		fs.EXPECT().DecodeFile(ctx, image).Return(mat, nil)
		sErr := errors.New("failed to show")
		id.EXPECT().Show(mat).Return(sErr)
		ps.EXPECT().SplitMatIntoTwoMats(ctx, mat).Return(expMatL, expMatR)
		uErr := errors.New("unMatch left image size")
		iv.EXPECT().ValidateImageSize(ctx, expMatL).Return(uErr)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{
				PreprocessingService: ps,
				ImageValidator:       iv,
				FileService:          fs,
				Debug:                true,
				ImageDebugger:        id,
			},
			d.SlamField{},
		)

		actlByte, actrByte, actStr, actErr := h.PreProcessStereoMerged(ctx, image)
		assert.Equal(t, uErr, actErr)
		assert.Nil(t, actlByte)
		assert.Nil(t, actrByte)
		assert.Empty(t, actStr)
	})
}

func Test_ProcessSlamMono(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ss := mock.NewMockSlamServiceInterface(ctrl)

	ctx := context.Background()
	logger := pkg.InitLogger(config.LogSettingsFilePath, config.ServiceName)
	ctx = pkg.WithLogger(ctx, logger)

	image := []byte("unused image")

	t.Run("ok_processslam_mono", func(t *testing.T) {
		pose := &kdslam.Pose{
			Tx: 2.2345678,
			Ty: 2.3345678,
			Tz: 2.4345678,
			Qx: 2.4345678,
			Qy: 2.5345678,
			Qz: 2.6345678,
			Qw: 2.7345678,
		}
		status := model.SlamState(gi.SlamState_MAP_LOADING_IN_PROGRESS)
		ss.EXPECT().GetPoseMono(ctx, image).Return(pose, status, nil)

		expected := createPubPose(pose, status)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{},
			d.SlamField{
				SlamService: ss,
			},
		)

		actual, err := h.ProcessSlamMono(ctx, image)
		assert.Equal(t, nil, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("ng_processslam_mono", func(t *testing.T) {

		err := errors.New("ng_processslam_mono")
		status := model.SlamState(gi.SlamState_TRACKING_LOST)
		ss.EXPECT().GetPoseMono(ctx, image).Return(nil, status, err)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{},
			d.SlamField{
				SlamService: ss,
			},
		)

		actual, acterr := h.ProcessSlamMono(ctx, image)
		assert.Error(t, err, acterr)
		assert.Empty(t, actual)
	})
}

func Test_ProcessSlamStereo(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ss := mock.NewMockSlamServiceInterface(ctrl)

	ctx := context.Background()
	logger := pkg.InitLogger(config.LogSettingsFilePath, config.ServiceName)
	ctx = pkg.WithLogger(ctx, logger)

	image := []byte("unused image")

	t.Run("ok_processslam_stereo", func(t *testing.T) {
		pose := &kdslam.Pose{
			Tx: 2.2345678,
			Ty: 2.3345678,
			Tz: 2.4345678,
			Qx: 2.4345678,
			Qy: 2.5345678,
			Qz: 2.6345678,
			Qw: 2.7345678,
		}
		status := model.SlamState(gi.SlamState_MAP_LOADING_IN_PROGRESS)
		ss.EXPECT().GetPoseStereo(ctx, image, image).Return(pose, status, nil)

		expected := createPubPose(pose, status)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{},
			d.SlamField{
				SlamService: ss,
			},
		)

		actual, err := h.ProcessSlamStereo(ctx, image, image)
		assert.Equal(t, nil, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("ng_processslam_stereo", func(t *testing.T) {

		err := errors.New("ng_processslam_stereo")
		status := model.SlamState(gi.SlamState_TRACKING_LOST)
		ss.EXPECT().GetPoseStereo(ctx, image, image).Return(nil, status, err)

		h := NewHandler(
			d.GrpcConnectorField{},
			d.PreprocessField{},
			d.SlamField{
				SlamService: ss,
			},
		)

		actual, acterr := h.ProcessSlamStereo(ctx, image, image)
		assert.Error(t, err, acterr)
		assert.Empty(t, actual)
	})
}
