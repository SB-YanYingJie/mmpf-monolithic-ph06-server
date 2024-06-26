package internal

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateSendMatRequest(t *testing.T) {
	wd, _ := os.Getwd()

	t.Run("ok: successful call returns err == nil", func(t *testing.T) {
		testValidMatsize(t, []int{640, 400}, filepath.Join(wd, "test_assets", "640_400.png"))
		testValidMatsize(t, []int{225, 225}, filepath.Join(wd, "test_assets", "225_225.jpeg"))
	})
	t.Run("ng: failed call returns err != nil", func(t *testing.T) {
		testInvalidMatsize(t, []int{639, 400}, filepath.Join(wd, "test_assets", "640_400.png"), "invalid image width. expected=639, actual=640")
		testInvalidMatsize(t, []int{640, 401}, filepath.Join(wd, "test_assets", "640_400.png"), "invalid image height. expected=401, actual=400")
		testInvalidMatsize(t, []int{300, 225}, filepath.Join(wd, "test_assets", "225_225.jpeg"), "invalid image width. expected=300, actual=225")
		testInvalidMatsize(t, []int{225, 300}, filepath.Join(wd, "test_assets", "225_225.jpeg"), "invalid image height. expected=300, actual=225")
	})
}

func testValidMatsize(t *testing.T, filesize []int, path string) {
	ctx := context.Background()
	validator := NewImageValidator(filesize[0], filesize[1])

	mat := openFileAsMatAsGrayScale(path)
	err := validator.ValidateImageSize(ctx, &mat)
	assert.NoError(t, err)
}
func testInvalidMatsize(t *testing.T, filesize []int, path, errString string) {
	ctx := context.Background()
	validator := NewImageValidator(filesize[0], filesize[1])

	mat := openFileAsMatAsGrayScale(path)
	err := validator.ValidateImageSize(ctx, &mat)
	assert.EqualError(t, err, errString)
}
