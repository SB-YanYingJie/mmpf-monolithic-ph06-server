package pkg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/xerrors"
)

type Entries struct {
	entries []zapcore.Entry
}

func (es *Entries) Hook(e zapcore.Entry) error {
	es.entries = append(es.entries, e)
	return nil
}
func (es *Entries) Init() {
	es.entries = []zapcore.Entry{}
}

// please execute test with verbose output.
// go test -v .
func TestLogger(t *testing.T) {
	ctx := context.Background()
	es := &Entries{}
	logger, _ := zap.NewDevelopment(zap.Hooks(es.Hook))
	ctx = WithLogger(ctx, logger)

	t.Run("string output test", func(t *testing.T) {
		es.Init()
		logger := GetLogger(ctx)
		logger.Info("test logging")
		assert.Equal(t, 1, len(es.entries))
		assert.Equal(t, "test logging", es.entries[0].Message)
		assert.Equal(t, zap.InfoLevel, es.entries[0].Level)
	})
	t.Run("error output test", func(t *testing.T) {
		es.Init()
		logger := GetLogger(ctx)
		err := xerrors.New("test error")
		// logging with error stacktrace.
		logger.Error("test logging2", zap.Error(err))
		assert.Equal(t, 1, len(es.entries))
		assert.Equal(t, "test logging2", es.entries[0].Message)
		assert.Equal(t, zap.ErrorLevel, es.entries[0].Level)
	})
}

func TestInitLogger(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		logger := InitLogger("./config.yml", "sevicename")
		assert.NotEmpty(t, logger)
	})
}
