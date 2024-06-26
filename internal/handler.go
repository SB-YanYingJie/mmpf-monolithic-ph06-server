package internal

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/KudanJP/KdSlamGo/kdslam"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	gi "github.com/machinemapplatform/grpc-interface/golang"
	"github.com/machinemapplatform/library/logging"
	"github.com/machinemapplatform/library/middleware"
	"github.com/machinemapplatform/library/model"
	"github.com/machinemapplatform/library/mytime"
	"github.com/machinemapplatform/mmpf-monolithic/cmd/config"
	d "github.com/machinemapplatform/mmpf-monolithic/internal/domain"
	im "github.com/machinemapplatform/mmpf-monolithic/internal/model"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	IMAGETYPE_MONO             = "mono"
	IMAGETYPE_STEREO_SEPARATED = "stereo_separated"
	IMAGETYPE_STEREO_MERGED    = "stereo_merged"
)

const (
	NUMBER_OF_LENSES_MONO   = "mono"
	NUMBER_OF_LENSES_STEREO = "stereo"
)

const (
	UnknownHost = "UNKNOWN_HOST"
)

type Handler struct {
	gi.MmpfMonolithicServer
	gm               d.GrpcConnectorField
	pm               d.PreprocessField
	sm               d.SlamField
	isSlamProcessing bool
	mutex            sync.Mutex
}

func NewHandler(gm d.GrpcConnectorField, pm d.PreprocessField, sm d.SlamField) Handler {
	return Handler{
		gm:               gm,
		pm:               pm,
		sm:               sm,
		isSlamProcessing: false,
	}
}

func (h *Handler) Slam(ctx context.Context, rawRequest *gi.SlamRequest) (*gi.SlamResponse, error) {
	startTime := mytime.NowUnixNano(ctx)
	logger := ctxzap.Extract(ctx)
	ctx = logging.WithLogger(ctx, logger)

	logger.Debug("SlamRequest",
		zap.Any("Metadata", rawRequest.Metadata),
		zap.Any("NumberOfLenses", rawRequest.NumberOfLenses),
		zap.Any("RequestTime", rawRequest.RequestTime),
	)

	h.mutex.Lock()
	if h.isSlamProcessing {
		logger.Debug("Slam process is disturbed")
		h.mutex.Unlock()
		return &gi.SlamResponse{}, nil
	}
	h.isSlamProcessing = true
	h.mutex.Unlock()

	defer func() {
		h.isSlamProcessing = false
	}()

	ctx = h.initCtxMetadata(ctx, rawRequest.Metadata, rawRequest.RequestTime)

	_, err := middleware.GetMetadata(ctx, model.MD_KEY_REQUEST_ID)
	if err != nil {
		logger.Debug("request_id not found")
	}

	request, err := h.gm.GrpcValidator.ValidateSlamRequest(ctx, rawRequest)
	if err != nil {
		logger.Warn("invalid argument", zap.Error(err))
		return getErrorResponse(ctx, rawRequest, model.INVALID_ARGUMENT, err.Error())
	}

	var result model.Pose
	var status model.ErrorStatus
	var path string

	switch config.ImageType {
	case IMAGETYPE_MONO:
		result, status, path, err = h.GetResultMono(ctx, request.Images[model.CENTER])
	case IMAGETYPE_STEREO_SEPARATED:
		result, status, path, err = h.GetResultStereoSeparated(ctx, request.Images[model.LEFT], request.Images[model.RIGHT])
	case IMAGETYPE_STEREO_MERGED:
		result, status, path, err = h.GetResultStereoMerged(ctx, request.Images[model.MERGED])
	default:
		panic("invalid env imageType")
	}
	if err != nil {
		logger.Error("failed to get slam result", zap.Error(err))
		return getErrorResponse(ctx, rawRequest, model.INTERNAL, err.Error())
	}

	ctx = middleware.WithMetadata(ctx, model.MD_KEY_RAW_IMAGE, path)
	err = h.PublishPose(ctx, result, status)
	if err != nil {
		logger.Warn("failed to publish pose", zap.Error(err))
	}

	metadata := safeMergeMap(middleware.Metadata(ctx), rawRequest.Metadata, "client_")
	response, err := toResponse(ctx, metadata, &result)

	endTime := mytime.NowUnixNano(ctx)
	elapsedTime := float64(endTime-startTime) / 1000000                                         // ナノ秒→ミリ秒
	response.Metadata[model.MD_KEY_ELAPSED_TIME] = strconv.FormatFloat(elapsedTime, 'f', 3, 64) // 少数第4位で四捨五入
	logger.Debug("metadata", zap.Any("md", response.Metadata))
	return response, err
}

func (h *Handler) initCtxMetadata(ctx context.Context, metadata map[string]string, t *timestamppb.Timestamp) context.Context {
	ctx = middleware.WithMetadata(ctx, model.MD_KEY_MMID, config.MMID)
	if t != nil {
		nanoT := strconv.FormatInt(t.AsTime().UnixNano(), 10)
		ctx = middleware.WithMetadata(ctx, model.MD_KEY_REQUEST_TIME, nanoT)
	}
	for k, v := range metadata {
		ctx = middleware.WithMetadata(ctx, k, v)
	}
	ctx = middleware.WithMetadata(ctx, model.MD_KEY_RAW_IMAGE_EXT, config.ImageExt)
	ctx = middleware.WithMetadata(ctx, model.MD_KEY_FRAME_SIZE_WIDTH, h.pm.FrameSizeWidth)
	ctx = middleware.WithMetadata(ctx, model.MD_KEY_FRAME_SIZE_HEIGHT, h.pm.FrameSizeHeight)
	ctx = middleware.WithMetadata(ctx, model.MD_KEY_REDIS_IMAGE_DB, h.pm.ImageDbNumber)
	ctx = middleware.WithMetadata(ctx, model.MD_KEY_TARGET_FPS, h.sm.FpsStr)
	return ctx
}

func (h *Handler) GetResultMono(ctx context.Context, rawImage []byte) (model.Pose, model.ErrorStatus, string, error) {
	logger := logging.GetLogger(ctx)

	image, path, err := h.PreProcessMono(ctx, rawImage)
	if err != nil {
		logger.Error("failed to preprocess mono", zap.Error(err))
		return model.Pose{}, model.INTERNAL, "", err
	}

	result, err := h.ProcessSlamMono(ctx, image)
	if err != nil {
		logger.Error("failed to process slam mono", zap.Error(err))
		return model.Pose{}, model.INTERNAL, "", err
	}

	return result, model.NO_ERROR, path, err
}

func (h *Handler) GetResultStereoSeparated(ctx context.Context, lRawImage, rRawImage []byte) (model.Pose, model.ErrorStatus, string, error) {
	logger := logging.GetLogger(ctx)

	lImage, rImage, lPath, err := h.PreProcessStereoSeparated(ctx, lRawImage, rRawImage)
	if err != nil {
		logger.Error("failed to preprocess stereo separated", zap.Error(err))
		return model.Pose{}, model.INTERNAL, "", err
	}

	result, err := h.ProcessSlamStereo(ctx, lImage, rImage)
	if err != nil {
		logger.Error("failed to process slam stereo", zap.Error(err))
		return model.Pose{}, model.INTERNAL, "", err
	}

	return result, model.NO_ERROR, lPath, err
}

func (h *Handler) GetResultStereoMerged(ctx context.Context, rawImage []byte) (model.Pose, model.ErrorStatus, string, error) {
	logger := logging.GetLogger(ctx)

	lImage, rImage, path, err := h.PreProcessStereoMerged(ctx, rawImage)
	if err != nil {
		logger.Error("failed to preprocess stereo merged", zap.Error(err))
		return model.Pose{}, model.INTERNAL, "", err
	}

	result, err := h.ProcessSlamStereo(ctx, lImage, rImage)
	if err != nil {
		logger.Error("failed to process slam stereo", zap.Error(err))
		return model.Pose{}, model.INTERNAL, "", err
	}

	return result, model.NO_ERROR, path, err
}

func (h *Handler) PreProcessMono(ctx context.Context, image []byte) ([]byte, string, error) {
	logger := logging.GetLogger(ctx)
	logger.Debug("start handle mono")

	// viewerで表示するためredisに生の画像を保存しておく
	path, err := h.pm.FileService.WriteFile(ctx, image)
	if err != nil {
		logger.Warn("failed to write file", zap.Error(err))
	}

	// kdslamに読み込ませるため画像をグレースケール化しておく
	mat, err := h.pm.FileService.DecodeFile(ctx, image)
	if err != nil {
		logger.Error("failed to decode file", zap.Error(err))
		return nil, "", err
	}
	if h.pm.Debug {
		err = h.pm.ImageDebugger.Show(mat)
		if err != nil {
			logger.Warn("failed to show debug image", zap.Error(err))
		}
	}
	defer mat.Close()

	if err := h.pm.ImageValidator.ValidateImageSize(ctx, mat); err != nil {
		logger.Error("unMatch image size", zap.Error(err))
		return nil, "", err
	}

	byteImage := h.pm.PreprocessingService.ConvertMatToBytes(ctx, mat)
	return byteImage, path, err
}

func (h *Handler) PreProcessStereoSeparated(ctx context.Context, lImage, rImage []byte) ([]byte, []byte, string, error) {
	logger := logging.GetLogger(ctx)
	logger.Debug("start handle stereo_separated")

	// viewerで表示するためredisに生の画像を保存しておく
	lPath, err := h.pm.FileService.WriteFile(ctx, lImage)
	if err != nil {
		logger.Warn("failed to write file", zap.Error(err))
	}

	// kdslamに読み込ませるため画像をグレースケール化しておく
	lMat, rMat, err := h.pm.FileService.DecodeFiles(ctx, lImage, rImage)
	if err != nil {
		logger.Error("failed to decode file", zap.Error(err))
		return nil, nil, "", err
	}
	if h.pm.Debug {
		err = h.pm.ImageDebugger.Show(lMat, rMat)
		if err != nil {
			logger.Warn("failed to show debug image", zap.Error(err))
		}
	}
	defer lMat.Close()
	defer rMat.Close()

	if err := h.pm.ImageValidator.ValidateImageSize(ctx, lMat); err != nil {
		logger.Error("unMatch left image size", zap.Error(err))
		return nil, nil, "", err
	}

	if err := h.pm.ImageValidator.ValidateImageSize(ctx, rMat); err != nil {
		logger.Error("unMatch right image size", zap.Error(err))
		return nil, nil, "", err
	}

	lByteImage := h.pm.PreprocessingService.ConvertMatToBytes(ctx, lMat)
	rByteImage := h.pm.PreprocessingService.ConvertMatToBytes(ctx, rMat)
	return lByteImage, rByteImage, lPath, err
}

func (h *Handler) PreProcessStereoMerged(ctx context.Context, image []byte) ([]byte, []byte, string, error) {
	logger := logging.GetLogger(ctx)
	logger.Debug("start handle stereo_merged")

	// viewerで表示するためredisに生の画像を保存しておく
	path, err := h.pm.FileService.WriteFile(ctx, image)
	if err != nil {
		logger.Warn("failed to write file", zap.Error(err))
	}

	// kdslamに読み込ませるため画像をグレースケール化しておく
	mat, err := h.pm.FileService.DecodeFile(ctx, image)
	if err != nil {
		logger.Error("failed to decode file", zap.Error(err))
		return nil, nil, "", err
	}
	if h.pm.Debug {
		err = h.pm.ImageDebugger.Show(mat)
		if err != nil {
			logger.Warn("failed to show debug image", zap.Error(err))
		}
	}
	defer mat.Close()

	lMat, rMat := h.pm.PreprocessingService.SplitMatIntoTwoMats(ctx, mat)
	defer lMat.Close()
	defer rMat.Close()

	if err := h.pm.ImageValidator.ValidateImageSize(ctx, lMat); err != nil {
		logger.Error("unMatch left image size", zap.Error(err))
		return nil, nil, "", err
	}
	if err := h.pm.ImageValidator.ValidateImageSize(ctx, rMat); err != nil {
		logger.Error("unMatch right image size", zap.Error(err))
		return nil, nil, "", err
	}

	lByteImage := h.pm.PreprocessingService.ConvertMatToBytes(ctx, lMat)
	rByteImage := h.pm.PreprocessingService.ConvertMatToBytes(ctx, rMat)
	return lByteImage, rByteImage, path, err
}

func (h *Handler) ProcessSlamMono(ctx context.Context, image []byte) (model.Pose, error) {
	logger := logging.GetLogger(ctx)
	logger.Debug("start slam mono")

	pose, state, err := h.sm.SlamService.GetPoseMono(ctx, image)
	if err != nil {
		logger.Warn("failed to GetPoseMono", zap.Error(err))
		return model.Pose{}, err
	}
	return createPubPose(pose, model.SlamState(state)), err
}

func (h *Handler) ProcessSlamStereo(ctx context.Context, lImage, rImage []byte) (model.Pose, error) {
	logger := logging.GetLogger(ctx)
	logger.Debug("start slam stereo")

	pose, state, err := h.sm.SlamService.GetPoseStereo(ctx, lImage, rImage)
	if err != nil {
		logger.Warn("failed to GetPoseStereo", zap.Error(err))
		return model.Pose{}, err
	}
	return createPubPose(pose, model.SlamState(state)), err
}

func createPubPose(pose *kdslam.Pose, state model.SlamState) model.Pose {
	return model.Pose{
		PosX:      pose.Tx,
		PosY:      pose.Ty,
		PosZ:      pose.Tz,
		QuatX:     pose.Qx,
		QuatY:     pose.Qy,
		QuatZ:     pose.Qz,
		QuatW:     pose.Qw,
		SlamState: model.SlamState(state),
	}
}

func (h *Handler) PublishPose(ctx context.Context, result interface{}, status model.ErrorStatus) error {
	logger := logging.GetLogger(ctx)
	if result == nil && status == model.NO_ERROR {
		return nil
	}
	j, err := json.Marshal(result)
	if err != nil {
		logger.Error("failed to marshal", zap.Error(err))
		return err
	}
	j2, err := middleware.ToDataWithContext(ctx, j)
	if err != nil {
		logger.Error("failed to make dataWithContext", zap.Error(err))
		return err
	}
	_, err = h.sm.Redis.Publish(ctx, config.RedisPubsubPoseChannel, j2)
	if err != nil {
		logger.Error("failed to publish", zap.Error(err))
		return err
	}
	return nil
}

func getErrorResponse(ctx context.Context, rawRequest *gi.SlamRequest, status model.ErrorStatus, message string) (*gi.SlamResponse, error) {
	logger := logging.GetLogger(ctx)
	hostname, herr := os.Hostname()
	if herr != nil {
		logger.Error("failed to get hostname")
		hostname = UnknownHost
	}
	t := mytime.NowUnixNano(ctx)
	return &gi.SlamResponse{
		Metadata: rawRequest.Metadata,
		Result:   toErrorResult(status, hostname, message, t),
	}, nil
}

func toErrorResult(status model.ErrorStatus, from, message string, t int64) *gi.SlamResponse_Error {
	return &gi.SlamResponse_Error{Error: &gi.ErrorResult{
		Status:    im.ToGrpcErrorStatus[status],
		Message:   message,
		From:      from,
		ErrorTime: toTimestamppb(t),
	},
	}
}

func toTimestamppb(t int64) *timestamppb.Timestamp {
	return timestamppb.New(time.Unix(0, t))
}

func safeMergeMap(map1, map2 map[string]string, prefix string) map[string]string {
	result := map[string]string{}
	for k, v := range map1 {
		result[k] = v
	}

	for k, v := range map2 {
		key := prefix + k
		if _, ok := result[key]; !ok {
			result[key] = v
		}
	}
	return result
}

func toResponse(ctx context.Context, metadata map[string]string, result *model.Pose) (*gi.SlamResponse, error) {
	logger := logging.GetLogger(ctx)

	tString := metadata[model.MD_KEY_REQUEST_TIME]
	t, err := strconv.ParseInt(tString, 10, 64)
	if err != nil {
		logger.Warn("t in ctx is NOT int type or undefined")
	}
	response := &gi.SlamResponse{
		Metadata: metadata,
	}
	response.Result = &gi.SlamResponse_Pose{
		Pose: &gi.PoseResult{
			RequestTime: toTimestamppb(t),
			PosX:        result.PosX,
			PosY:        result.PosY,
			PosZ:        result.PosZ,
			QuatX:       result.QuatX,
			QuatY:       result.QuatY,
			QuatZ:       result.QuatZ,
			QuatW:       result.QuatW,
			SlamState:   im.ToGrpcSlamState[result.SlamState],
		},
	}
	logger.Debug("SlamResponse", zap.Any("response", response))
	return response, nil
}
