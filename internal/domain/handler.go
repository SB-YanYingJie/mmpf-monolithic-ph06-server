package domain

import (
	"github.com/machinemapplatform/library/file"
	"github.com/machinemapplatform/library/redis"
)

type GrpcConnectorField struct {
	GrpcValidator GrpcValidatorInterface
}

type PreprocessField struct {
	FileService          FileServiceInterface
	ImageValidator       ImageValidatorInterface
	PreprocessingService PreprocessingServiceInterface
	NumberOfLenses       string
	Debug                bool
	ImageDebugger        ImageDebuggerInterface
	FrameSizeWidth       string
	FrameSizeHeight      string
	ImageDbNumber        string
	File                 file.FileInterface
}

type SlamField struct {
	SlamService SlamServiceInterface
	FpsStr      string
	Redis       redis.RedisInterface
}
