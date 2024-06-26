package internal

import (
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gocv.io/x/gocv"
)

func Test_trim(t *testing.T) {

	inputPath := "./test_assets/1280_400.png"
	format := "png"
	mat := openFileAsMatAsGrayScale(inputPath)
	lMat, rMat := trim(&mat, 0, 0, 640, 400)
	lImage, _ := lMat.ToImage()
	rImage, _ := rMat.ToImage()

	assert.Equal(t, 640, lImage.Bounds().Dx())
	assert.Equal(t, 400, lImage.Bounds().Dy())
	assert.Equal(t, 640, rImage.Bounds().Dx())
	assert.Equal(t, 400, rImage.Bounds().Dy())

	// Trim後のファイル確認用
	switch format {
	case "png":
		lFile, _ := os.Create("./test_assets/left.png")
		defer lFile.Close()
		if err := png.Encode(lFile, lImage); err != nil {
			t.Fatal("Could not write left file")
		}
		rFile, _ := os.Create("./test_assets/right.png")
		defer rFile.Close()
		if err := png.Encode(rFile, lImage); err != nil {
			t.Fatal("Could not write left file")
		}
	case "jpeg":
		lFile, _ := os.Create("./test_assets/left.jpeg")
		defer lFile.Close()
		if err := jpeg.Encode(lFile, lImage, &jpeg.Options{Quality: 100}); err != nil {
			t.Fatal("Could not write left file")
		}
		rFile, _ := os.Create("./test_assets/right.jpeg")
		defer rFile.Close()
		if err := jpeg.Encode(rFile, lImage, &jpeg.Options{Quality: 100}); err != nil {
			t.Fatal("Could not write left file")
		}
	}
}

func openFileAsMatAsGrayScale(path string) gocv.Mat {
	decoded := gocv.IMRead(path, gocv.IMReadGrayScale)
	return decoded
}
