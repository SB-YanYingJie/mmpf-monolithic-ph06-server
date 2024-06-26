package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	mock "github.com/machinemapplatform/library/file/mock"
	"github.com/machinemapplatform/library/middleware"
	"github.com/machinemapplatform/library/model"
	"github.com/machinemapplatform/mmpf-monolithic/cmd/config"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func Test_DecodeFile(t *testing.T) {
	ctx := context.Background()
	data, err := os.ReadFile("test_assets/225_225.jpeg")
	if err != nil {
		t.Error("image file not found")
		t.FailNow()
	}

	t.Run("ok: successful call returns err == nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		file := mock.NewMockFileInterface(ctrl)
		service := NewFileService(file, "mmid")

		mat, err := service.DecodeFile(ctx, data)
		assert.NotEmpty(t, mat)
		assert.NoError(t, err)
	})
}
func Test_WriteFile(t *testing.T) {
	ctx := context.Background()

	tStr := "1626437594000000015"
	mmid := "mmid"

	t.Run("ng: failed call returns err != nil(failed get t)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		file := mock.NewMockFileInterface(ctrl)
		service := NewFileService(file, mmid)
		ctx = middleware.WithMetadata(ctx, model.MD_KEY_REQUEST_TIME, "aaa")

		path, err := service.WriteFile(ctx, []byte("image data"))
		assert.Empty(t, path)
		assert.Error(t, err)
	})
	t.Run("ng: failed call returns err != nil(failed store)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		file := mock.NewMockFileInterface(ctrl)
		service := NewFileService(file, mmid)
		ctx = middleware.WithMetadata(ctx, model.MD_KEY_REQUEST_TIME, tStr)
		file.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("failed to store"))

		path, err := service.WriteFile(ctx, []byte("image data"))
		assert.Empty(t, path)
		assert.Error(t, err)
	})
	t.Run("ok: successful call returns err == nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		file := mock.NewMockFileInterface(ctrl)
		service := NewFileService(file, mmid)
		ctx = middleware.WithMetadata(ctx, model.MD_KEY_REQUEST_TIME, tStr)

		// expectedDir := filepath.Join("2021_07_16/12/mmid/image")
		expectedFilename := strings.Join([]string{tStr, "c", config.ImageExt}, ".")
		expectedPath := filepath.Join("2021_07_16/12/mmid/raw_image", expectedFilename)
		file.EXPECT().Store(gomock.Any(), gomock.Any(), expectedPath).Return(nil)

		path, err := service.WriteFile(ctx, []byte("image data"))
		assert.Equal(t, expectedPath, path)
		assert.NoError(t, err)
	})
}

func Test_DecodeFiles(t *testing.T) {
	ctx := context.Background()

	t.Run("ok: successful call returns err == nil", func(t *testing.T) {
		data, err := os.ReadFile("test_assets/225_225.jpeg")
		if err != nil {
			t.Error("image file not found")
			t.FailNow()
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		file := mock.NewMockFileInterface(ctrl)
		service := NewFileService(file, "mmid")

		lImage, rImage, err := service.DecodeFiles(ctx, data, data)
		assert.NotEmpty(t, lImage)
		assert.NotEmpty(t, rImage)
		assert.NoError(t, err)
	})
}
