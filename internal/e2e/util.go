package e2e

import (
	"io/ioutil"
	"log"
	"os"
)

const address = "localhost:50051"

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
