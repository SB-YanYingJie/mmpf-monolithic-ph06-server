package e2e

import (
	"context"
	"log"

	gi "github.com/machinemapplatform/grpc-interface/golang"
	"github.com/machinemapplatform/library/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DeviceClient struct {
	C    gi.MmpfMonolithicClient
	conn *grpc.ClientConn
}

func NewDeviceClient(address string) *DeviceClient {
	c, conn := initDeviceClientCore(address)
	return &DeviceClient{C: c, conn: conn}
}

func initDeviceClientCore(address string) (gi.MmpfMonolithicClient, *grpc.ClientConn) {
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

func (d *DeviceClient) Slam(ctx context.Context, images []*gi.Image, numberOfLenses gi.NumberOfLenses, requestId string, t *timestamppb.Timestamp) (*gi.SlamResponse, error) {
	metadata := map[string]string{model.MD_KEY_REQUEST_ID: requestId}
	log.Printf("d.C.SlamReq: %s", requestId)

	res, err := d.C.Slam(ctx, &gi.SlamRequest{
		Metadata:       metadata,
		RequestTime:    t,
		NumberOfLenses: numberOfLenses,
		Images:         images,
	})
	log.Printf("d.C.SlamRes: %s", requestId)
	if err != nil {
		log.Printf("SlamErr: %s\n", err)
	}

	return res, err
}

func (c *DeviceClient) Close() {
	c.conn.Close()
}
