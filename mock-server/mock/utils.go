package mock

import (
	"io/ioutil"
	"log"
	"os"
)

type FileOpener interface {
	OpenFileAsBytes(filePath string) []byte
}

type ApitestFileOpener struct{}

func (fo *ApitestFileOpener) OpenFileAsBytes(filePath string) []byte {
	return []byte(filePath)
}

func Chunk(rawBytes []byte, chunkSize int) [][]byte {
	rawBytesLen := len(rawBytes)
	chunked := [][]byte{}
	for i := 0; i < rawBytesLen; i += chunkSize {
		end := i + chunkSize
		if rawBytesLen < end {
			end = rawBytesLen
		}
		chunked = append(chunked, rawBytes[i:end])
	}
	return chunked
}

func OpenFileAsBytes(filePath string) []byte {
	// ファイルを開く
	file, opFileErr := os.Open(filePath)
	if opFileErr != nil {
		log.Printf("file open error: %v", opFileErr)
	}
	defer file.Close()
	// byte列に変換
	b, readFileErr := ioutil.ReadAll(file)
	if readFileErr != nil {
		log.Printf("file read error: %v", readFileErr)
	}
	return b
}
