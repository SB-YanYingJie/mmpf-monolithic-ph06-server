//go:build xWindow
// +build xWindow

package internal

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/machinemapplatform/mmpf-monolithic/cmd/config"
	"gocv.io/x/gocv"
)

//以下のコードは、slamが動作することをX Windowで確認するため残しています。
//ローカルでの実行を想定しており、ci環境ではまわさないので常時は実行しません。
//-tags=xWindow オプションを付けた場合のみ、テスト実行可能です。
func TestProcessFrame(t *testing.T) {

	slam := NewSlam(log.New(os.Stdout, "test", 2))
	if err := slam.StartSlam(config.CalibPath, config.VocabPath, config.Fps); err != nil {
		fmt.Printf("failed to StartSlam:%s", err)
	}
	slam.SetAutoExpansion(false)
	if err := slam.LoadMap(config.KdmpPath, 0, 0); err != nil {
		log.Panicf("faile to LoadMap:%s", err)
	}
	w, h, _ := slam.GetImageSize()
	fmt.Printf("w : %d", w)
	fmt.Printf("h : %d", h)
	debugWindow := gocv.NewWindow("debug window1")

	t.Run("stereo_ok", func(t *testing.T) {
		imagebytesLeft, _ := ioutil.ReadFile(EditImageFile("/app/internal/test_assets/left.png"))
		imagebytesRight, _ := ioutil.ReadFile(EditImageFile("/app/internal/test_assets/right.png"))
		for {
			debug, err := slam.ProcessFrameStereo(imagebytesLeft, imagebytesRight, true)
			if err != nil {
				fmt.Printf("debug:%s", err)
			}
			fmt.Println(len(debug))
			debugMat, _ := gocv.NewMatFromBytes(h, w, gocv.MatTypeCV8UC4, debug)
			if err != nil {
				fmt.Printf("mat:%s", err)
			}
			debugWindow.IMShow(debugMat)
			debugWindow.WaitKey(1)
		}
	})

	t.Run("mono_ok", func(t *testing.T) {
		imagebytes, err := ioutil.ReadFile(EditImageFile("/app/internal/test_assets/1267_1.png"))
		fmt.Printf("len:%d\n", len(imagebytes))
		if err != nil {
			fmt.Printf("file:%s", err)
		}

		for {
			debug, err := slam.ProcessFrame(imagebytes, true)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(len(debug))
			debugMat, err := gocv.NewMatFromBytes(h, w, gocv.MatTypeCV8UC4, debug)
			if err != nil {
				fmt.Println(err.Error())
			}

			debugWindow.IMShow(debugMat)
			debugWindow.WaitKey(1)
		}
	})
}

func TestProcessFrameWithGoCVDecode(t *testing.T) {

	slam := NewSlam(log.New(os.Stdout, "test", 2))
	if err := slam.StartSlam(config.CalibPath, config.VocabPath, config.Fps); err != nil {
		fmt.Printf("failed to StartSlam:%s", err.Error())
	}
	slam.SetAutoExpansion(false)
	if err := slam.LoadMap(config.KdmpPath, 0, 0); err != nil {
		fmt.Printf("failed to StartSlam:%s", err.Error())
	}
	w, h, _ := slam.GetImageSize()
	fmt.Printf("w : %d", w)
	fmt.Printf("h : %d", h)
	debugWindow := gocv.NewWindow("debug window1")

	t.Run("mono_ok", func(t *testing.T) {
		imagebytes, err := ioutil.ReadFile(editImageFileWithGoCV("/app/internal/test_assets/1267_1.png"))
		fmt.Printf("len:%d\n", len(imagebytes))
		if err != nil {
			fmt.Printf("file:%s", err)
		}

		for {
			debug, err := slam.ProcessFrame(imagebytes, true)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(len(debug))
			debugMat, err := gocv.NewMatFromBytes(h, w, gocv.MatTypeCV8UC4, debug)
			if err != nil {
				fmt.Println(err.Error())
			}
			debugWindow.IMShow(debugMat)
			debugWindow.WaitKey(1)
		}
	})
}

func EditImageFile(path string) string {

	fOrigin, err := os.Open(path)
	if err != nil {
		fmt.Printf("edit1:%s", err)
	}

	img, err := png.Decode(fOrigin)
	if err != nil {
		fmt.Printf("edit2:%s", err)
	}
	mat, err := gocv.ImageToMatRGB(img)
	if err != nil {
		fmt.Printf("edit3:%s", err)
	}
	matGray := gocv.NewMat()
	gocv.CvtColor(mat, &matGray, gocv.ColorRGBToGray)
	imgb, err := matGray.DataPtrUint8()
	if err != nil {
		fmt.Printf("edit4:%s", err)
	}
	filePath := strings.Replace(path, ".png", "Edited.bin", -1)
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("edit5:%s", err)
	}
	defer f.Close()
	_, err = f.Write([]byte(imgb))
	if err != nil {
		fmt.Printf("failed to f.Write:%s", err)
	}
	fmt.Print(filePath)
	return filePath

}

func editImageFileWithGoCV(path string) string {

	// Use gocv.IMReadUnchanged for grayscale images,
	// gocv.IMReadGrayScale for color images
	mat := gocv.IMRead(path, gocv.IMReadGrayScale)
	imgb, err := mat.DataPtrUint8()
	if err != nil {
		fmt.Printf("edit4:%s", err)
	}
	filePath := strings.Replace(path, ".png", "Edited.bin", -1)
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("edit5:%s", err)
	}
	defer f.Close()
	_, err = f.Write([]byte(imgb))
	if err != nil {
		fmt.Printf("failed to f.Write:%s", err)
	}
	fmt.Print(filePath)
	return filePath

}
