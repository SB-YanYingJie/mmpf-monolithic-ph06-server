package config

import (
	"github.com/machinemapplatform/library/env"
)

// common
var (
	// LogSettingsFilePath is log setting file path
	LogSettingsFilePath = env.ParseString("LOG_SETTINGS_FILE_PATH", nil)
	// MMID is machinemap id of service
	MMID = env.ParseString("MMID", nil)
	// Port is listening port of service
	Port = env.ParseString("PORT", nil)
	// ServiceName is name of service
	ServiceName = env.ParseString("SERVICE_NAME_MONOLITHIC", env.String("service_name_not_specified"))
)

// grpc
var (
	// Creds is credential
	// Creds = env.ParseString("CREDENTIAL", nil)
	// ImageExt is file extension of image
	ImageExt = env.ParseString("SAVE_IMAGE_EXT", nil)
)

// preprocess
var (
	// DevDisplayRawImage is a flag determin display raw image or not
	DevDisplayRawImage = env.ParseBool("DEV_DISPLAY_RAW_IMAGE", env.Bool(false))
	// ImageType takes mono, stereo_merged and stereo_separated
	ImageType = env.ParseString("SEND_IMAGE_TYPE", nil)
	// NumberOfLenses is the number of lenses in the client
	NumberOfLenses = env.ParseString("NUMBER_OF_LENSES", nil)
	// TrimmingParameter is reference value or range for image segmentation
	TrimmingParameter = env.ParseIntArray("TRIMMING_PARAMETER", []int{0, 0, 0, 0})
)

// slam
var (
	// CalibPath is Calibration file Path
	CalibPath = env.ParseString("KD_CALIB_PATH", nil)
	// DevDisplayDebugImage is a flag determin display image or not
	DevDisplayDebugImage = env.ParseBool("DEV_DISPLAY_DEBUG_IMAGE", env.Bool(false))
	// Fps is frame rate
	Fps = env.ParseFloat32("TARGET_FPS", nil)
	// KdmpPath is kudan-map file Path
	KdmpPath = env.ParseString("KD_MAP_PATH", nil)
	// MapExpansionFlag is the flag to expansion a map
	MapExpansionFlag = env.ParseBool("KD_MAP_EXPAND", env.Bool(false))
	// MapId is map id service applies to
	MapId = env.ParseString("MAPID", nil)
	// VocabPath is vocabulary file Path
	VocabPath = env.ParseString("KD_VOCAB_PATH", nil)
)

// redis
var (
	// ImageStoreRedisTtl is Redis TTL if ImageStore is redis (default:3sec)
	ImageStoreRedisTtl = env.ParseInt("IMAGE_STORE_REDIS_TTL", env.Int(3))
	// RedisAddress is address of redis
	RedisAddress = env.ParseString("REDIS_ADDRESS", nil)
	// RedisIdleTimeoutSeconds is time in seconds to close an idle connection
	RedisIdleTimeoutSeconds = env.ParseInt("REDIS_IDLE_TIMEOUT_SECONDS", nil)
	// RedisMaxIdle is maximun number of idle connection
	RedisMaxIdle = env.ParseInt("REDIS_MAX_IDLE", nil)
	// RedisPubsubPoseChannel is redis pubusub pose channel
	RedisPubsubPoseChannel = env.ParseString("REDIS_PUBSUB_CHANNEL_POSE", nil)
	// RedisPubsubDb is redis pubsub db
	RedisPubsubDb = env.ParseInt("REDIS_PUBSUB_DB", nil)
)
