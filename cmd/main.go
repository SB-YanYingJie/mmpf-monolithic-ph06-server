package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	gi "github.com/machinemapplatform/grpc-interface/golang"
	"github.com/machinemapplatform/library/file"
	"github.com/machinemapplatform/library/logging"
	"github.com/machinemapplatform/library/redis"
	"github.com/machinemapplatform/mmpf-monolithic/cmd/config"
	"github.com/machinemapplatform/mmpf-monolithic/internal"
	d "github.com/machinemapplatform/mmpf-monolithic/internal/domain"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
)

const (
	IMAGETYPE_MONO             = "mono"
	IMAGETYPE_STEREO_SEPARATED = "stereo_separated"
	IMAGETYPE_STEREO_MERGED    = "stereo_merged"
)

func main() {
	logSettingPath, err := filepath.Abs(config.LogSettingsFilePath)
	if err != nil {
		panic(err)
	}
	logger := logging.InitLogger(logSettingPath, config.ServiceName, config.MMID)

	ctx := context.Background()
	ctx = logging.WithLogger(ctx, logger)
	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		switch e := recover().(type) {
		case nil:
		case error:
			logger.Panic("panic", zap.Error(e))
		default:
			logger.Panic("panic", zap.Any("error", e))
		}
	}()

	netInterfaceAddresses, _ := net.InterfaceAddrs()
	for _, netInterfaceAddress := range netInterfaceAddresses {
		networkIP, ok := netInterfaceAddress.(*net.IPNet)
		if ok && !networkIP.IP.IsLoopback() && networkIP.IP.To4() != nil {
			ip := networkIP.IP.String()
			fmt.Println("Resolved Host IP: " + ip)
		}
	}

	lis, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcValidator := internal.NewGrpcValidator()

	redis := redis.NewRedis(config.RedisAddress, config.RedisMaxIdle, config.RedisIdleTimeoutSeconds, config.RedisPubsubDb)
	filer := file.NewFileRedis(redis, config.ImageStoreRedisTtl)
	fileService := internal.NewFileService(filer, config.MMID)

	preprocessingService := internal.NewPreprocessingService(config.TrimmingParameter[0], config.TrimmingParameter[1], config.TrimmingParameter[2], config.TrimmingParameter[3])
	imageDbNumber := strconv.Itoa(config.RedisPubsubDb)

	var rawImageDebugWindowNames []string
	if config.ImageType == IMAGETYPE_STEREO_SEPARATED {
		rawImageDebugWindowNames = []string{"raw_image_" + config.MMID + "_l", "raw_image_" + config.MMID + "_r"}
	} else {
		rawImageDebugWindowNames = []string{"raw_image_" + config.MMID}
	}
	rawImageDebugger, err := internal.NewImageDebugger(config.DevDisplayRawImage, rawImageDebugWindowNames...)
	if err != nil {
		logger.Warn("failed to init rawImageDebugger", zap.Error(err))
	}

	slamedImageDebugger, err := internal.NewImageDebugger(config.DevDisplayDebugImage, "image_"+config.MMID)
	if err != nil {
		logger.Warn("failed to init slamedImageDebugger", zap.Error(err))
	}

	slamLogger, err := zap.NewStdLogAt(logger, zap.DebugLevel)
	if err != nil {
		logger.Error("NewStdLog error: %s", zap.Error(err))
	}

	slamservice := internal.NewSlamService(
		config.MMID,
		config.CalibPath,
		config.VocabPath,
		config.KdmpPath,
		config.Fps,
		config.DevDisplayDebugImage,
		config.MapExpansionFlag,
		internal.NewSlam(slamLogger),
		slamedImageDebugger,
	)

	frameSizeWidth, frameSizeHeight, err := slamservice.Start(ctx)
	if err != nil {
		logger.Error("could not start slamService", zap.Error(err))
	}
	defer slamservice.Close(ctx)

	imageValidator := internal.NewImageValidator(frameSizeWidth, frameSizeHeight)
	frameSizeWidthStr, frameSizeHeightStr := strconv.Itoa(frameSizeWidth), strconv.Itoa(frameSizeHeight)

	err = redis.BSet(ctx, config.MapId+":"+config.MMID, []byte(config.RedisPubsubPoseChannel))
	if err != nil {
		logger.Warn("failed to Set mmid to Redis DB", zap.Error(err))
	}
	ch := make(chan os.Signal, 10)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		cancel()
		close(ch)
		if err = redis.BSet(ctx, config.MapId+":"+config.MMID, []byte("")); err != nil {
			logger.Debug("failed to BSet.")
		}
		logger.Debug("stop slam service.")
	}()

	handler := internal.NewHandler(
		d.GrpcConnectorField{
			GrpcValidator: *grpcValidator,
		},
		d.PreprocessField{
			FileService:          fileService,
			ImageValidator:       imageValidator,
			PreprocessingService: preprocessingService,
			NumberOfLenses:       config.NumberOfLenses,
			Debug:                config.DevDisplayRawImage,
			ImageDebugger:        rawImageDebugger,
			FrameSizeWidth:       frameSizeWidthStr,
			FrameSizeHeight:      frameSizeHeightStr,
			ImageDbNumber:        imageDbNumber,
			File:                 filer,
		},
		d.SlamField{
			SlamService: slamservice,
			FpsStr:      fmt.Sprintf("%.2f", config.Fps),
			Redis:       redis,
		},
	)

	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(generateRecoveryFunc(logger)),
	}
	// zapOptでinfoLevelのコードをすべてDebugLevelに変更
	zapOpt := grpc_zap.WithLevels(
		func(c codes.Code) zapcore.Level {
			var l zapcore.Level
			switch c {
			case codes.OK:
				l = zapcore.DebugLevel
			case codes.Canceled:
				l = zapcore.DebugLevel
			case codes.InvalidArgument:
				l = zapcore.DebugLevel
			case codes.NotFound:
				l = zapcore.DebugLevel
			case codes.AlreadyExists:
				l = zapcore.DebugLevel
			case codes.Unauthenticated:
				l = zapcore.DebugLevel
			case codes.DeadlineExceeded:
				l = zapcore.DebugLevel
			case codes.PermissionDenied:
				l = zapcore.DebugLevel
			default:
				l = zapcore.DebugLevel
			}
			return l
		},
	)

	gsvr := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger, zapOpt),
			grpc_recovery.UnaryServerInterceptor(opts...),
		)),
	)
	gi.RegisterMmpfMonolithicServer(gsvr, &handler)
	if err := gsvr.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func generateRecoveryFunc(logger *zap.Logger) func(p interface{}) error {
	return func(p interface{}) error {
		switch e := p.(type) {
		case error:
			logger.Panic("panic gRPC", zap.Error(e))
		default:
			logger.Panic("panic gRPC", zap.Any("error", e))
		}
		return status.Errorf(codes.Internal, "Unexpected error")
	}
}
