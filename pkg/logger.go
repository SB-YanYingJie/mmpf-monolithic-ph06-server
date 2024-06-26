package pkg

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type loggerKey string

const key loggerKey = "LoggerKey"

// WithLogger set Logger to context.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, key, logger)
}

// GetLogger get Logger from context.
func GetLogger(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(key).(*zap.Logger); ok {
		return logger
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("zap.NewDevelopment() failed: %+w", err))
	}
	return logger
}

func InitLogger(pathToConfigFile, serviceName string) *zap.Logger {
	ymlFile, err := os.ReadFile(pathToConfigFile)
	if err != nil {
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		return addDefaultField(logger, serviceName)
	}
	var inputConfig zap.Config
	if err := yaml.Unmarshal(ymlFile, &inputConfig); err != nil {
		panic(err)
	}

	config := zap.NewProductionConfig()
	config.Level = inputConfig.Level

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return addDefaultField(logger, serviceName)
}

func addDefaultField(logger *zap.Logger, serviceName string) *zap.Logger {
	return logger.With(
		zap.String("ServiceName", serviceName),
	)
}
