package mock

import (
	"context"
	"fmt"
	"log"

	gi "github.com/machinemapplatform/grpc-interface/golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Sender interface {
	Slam(ctx context.Context, fo FileOpener, pathToFile string, limit int)
	Close()
}

type DeviceClient struct {
	c    gi.MmpfMonolithicClient
	conn *grpc.ClientConn
}

type MonoClient DeviceClient

func initClientCore(address string) (gi.MmpfMonolithicClient, *grpc.ClientConn) {

	//TLSをオフにすると、ベーシック認証を残すことができない。
	//https://github.com/machinemapplatform/mmpf-monolithic/issues/19
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return gi.NewMmpfMonolithicClient(conn), conn
}

func NewMonoClient(address string) *MonoClient {
	c, conn := initClientCore(address)
	return &MonoClient{c: c, conn: conn}
}

func (m *MonoClient) Close() {
	m.conn.Close()
}

func (m *MonoClient) Slam(ctx context.Context, fo FileOpener, pathToFile string, limit int) {
	imgByte := fo.OpenFileAsBytes(pathToFile + "1267_0.png")
	res, err := m.c.Slam(ctx, &gi.SlamRequest{
		Metadata:       map[string]string{"k1": "v1", "delayed": "100"},
		NumberOfLenses: gi.NumberOfLenses_MONO,
		Images: []*gi.Image{
			{
				LensPlacement: gi.LensPlacement_CENTER,
				Byte:          imgByte,
			},
		},
		RequestTime: timestamppb.Now(),
	})
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println("res", res)
}
