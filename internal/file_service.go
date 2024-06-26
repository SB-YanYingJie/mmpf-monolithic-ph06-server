package internal

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/machinemapplatform/library/file"
	"github.com/machinemapplatform/library/logging"
	"github.com/machinemapplatform/library/middleware"
	"github.com/machinemapplatform/library/model"
	"github.com/machinemapplatform/mmpf-monolithic/cmd/config"
	"go.uber.org/zap"
	"gocv.io/x/gocv"
	"golang.org/x/sync/errgroup"
)

const (
	timeFormat     = "2006_01_02/15"
	dirFormat      = "%s/%s/raw_image" // example: 2021_07_27/20/mmid01/image
	filenameFormat = "%s.%s.%s"        // example: 1626437594000000015.c.png
)

func NewFileService(
	f file.FileInterface,
	mmid string,
) *FileService {
	return &FileService{
		f:    f,
		mmid: mmid,
	}
}

type FileService struct {
	f    file.FileInterface
	mmid string
}

func (f *FileService) DecodeFile(ctx context.Context, data []byte) (*gocv.Mat, error) {
	mat, err := gocv.IMDecode(data, gocv.IMReadGrayScale)
	if err != nil {
		return nil, fmt.Errorf("gocv IMDecode failed: %+w", err)
	}
	if mat.Empty() {
		return nil, fmt.Errorf("IMDecode returns Empty")
	}
	return &mat, nil
}

func (f *FileService) DecodeFiles(ctx context.Context, lData, rData []byte) (lImage, rImage *gocv.Mat, err error) {
	logger := logging.GetLogger(ctx)

	eg := errgroup.Group{}

	eg.Go(func() error {
		l, err := gocv.IMDecode(lData, gocv.IMReadGrayScale)
		if err != nil {
			logger.Error("gocv IMDecode ldata failed", zap.Error(err))
			return err
		}
		lImage = &l
		return nil
	})
	eg.Go(func() error {
		r, err := gocv.IMDecode(rData, gocv.IMReadGrayScale)
		if err != nil {
			logger.Error("gocv IMDecode rdata failed", zap.Error(err))
			return err
		}
		rImage = &r
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, nil, fmt.Errorf("readFiles failed: %+w", err)
	}

	return lImage, rImage, nil
}

func (f *FileService) WriteFile(ctx context.Context, image []byte) (string, error) {
	logger := logging.GetLogger(ctx)

	t, err := middleware.GetMetadataAsInt64(ctx, model.MD_KEY_REQUEST_TIME)
	if err != nil {
		logger.Error("failed to get t from metadata")
		return "", err
	}
	dir := buildWriteDir(t, f.mmid)
	filename := buildFileName(t, model.CENTER, config.ImageExt)
	path := filepath.Join(dir, filename)
	if err := f.f.Store(ctx, image, path); err != nil {
		logger.Error("failed to store", zap.String("path", path))
		return "", fmt.Errorf("failed to store: %+w", err)
	}
	return path, nil
}

// buildWriteDir returns dir to store image
// example: /2021_07_27/20/mmid01/image
func buildWriteDir(t int64, mmid string) string {
	return fmt.Sprintf(dirFormat, time.Unix(0, t).Format(timeFormat), mmid)
}

// buildFileName returns filename
// example: 1626437594000000015.c.png
func buildFileName(t int64, camPosType model.LensPlacement, ext string) string {
	return fmt.Sprintf(filenameFormat, strconv.Itoa(int(t)), camPosType, ext)
}
