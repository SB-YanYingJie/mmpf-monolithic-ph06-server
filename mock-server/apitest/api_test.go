package apitest

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	mock "github.com/machinemapplatform/mmpf-monolithic/mock-server/mock"
)

const (
	// address = "localhost"
	address = "172.18.0.1" // please set docker-network-gateway-ip
	// address    = "52.194.188.111" // please set ECS public ip
	pathToFile = "../assets/"
	duration   = time.Second
)

type MockFileOpener struct{}

func (fo MockFileOpener) OpenFileAsBytes(filePath string) []byte {
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

func TestHello(t *testing.T) {
	t.Run("Slam Connect ok", func(t *testing.T) {
		fmt.Println("SlamRequest ok started")
		port := os.Getenv("PORT")
		address := address + ":" + port
		c := mock.NewMonoClient(address)
		defer c.Close()
		ctx := context.Background()
		fo := MockFileOpener{}
		for i := 0; i < 10; i++ {
			request(c, ctx, fo, pathToFile, 2)
		}
	})
}

func request(s mock.Sender, ctx context.Context, fo mock.FileOpener, pathToFile string, limit int) {
	s.Slam(ctx, fo, pathToFile, limit)
}
