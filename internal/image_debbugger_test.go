package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugger(t *testing.T) {
	// 以下のようなテストを書いたが CI で実行できないためコメントアウトする。
	// t.Run("ng: containing empty name returns err when debug on", func(t *testing.T) {
	// 	_, err := NewImageDebugger(true, "aaa", "")
	// 	assert.Error(t, err)
	// })
	// t.Run("ok: containing only nonempty name returns no err when debug on.", func(t *testing.T) {
	// 	debugger, err := NewImageDebugger(true, "aaa", "bbb")
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, 2, debugger.windowsLen)

	// successful Show cannot be tested cause of depending on gocv(and that must be visible)
	// err = debugger.Show(&gocv.Mat{})
	// assert.Error(t, err)
	// err = debugger.Show(&gocv.Mat{}, &gocv.Mat{}, &gocv.Mat{})
	// assert.Error(t, err)
	// })
	t.Run("ok: containing any name returns no err when debug off", func(t *testing.T) {
		debugger, err := NewImageDebugger(false, "", "")
		assert.NoError(t, err)
		assert.Equal(t, 0, debugger.windowsLen)
	})
}
