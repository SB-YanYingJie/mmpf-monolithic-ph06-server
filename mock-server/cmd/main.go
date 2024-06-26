package main

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	gimock "github.com/machinemapplatform/grpc-interface/golang"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	delayed = "delayed"
)

type server struct {
	gimock.MmpfMonolithicServer
}

func (s *server) Slam(ctx context.Context, request *gimock.SlamRequest) (*gimock.SlamResponse, error) {
	log.Printf("SlamRequest started")
	log.Printf("request.Metadata: %v\n", request.Metadata)
	log.Printf("request.NumberOfLenses: %v\n", request.NumberOfLenses)
	log.Printf("request.RequestTime: %v\n", request.RequestTime)

	delayTimeMS := request.Metadata[delayed]
	if delayTimeMS != "" {
		err := delay(delayTimeMS)
		if err != nil {
			return nil, err
		}
	}

	metadata := map[string]string{
		delayed: delayTimeMS,
	}

	metadata = setMetadata(request.Metadata, metadata)
	metadata["elapsed_time"] = "33"
	metadata["request_time"] = request.GetRequestTime().AsTime().String()
	metadata["number_of_lenses"] = request.NumberOfLenses.String()
	metadata = setImages(request.GetImages(), metadata)

	return &gimock.SlamResponse{
		Metadata: metadata,
		Result: &gimock.SlamResponse_Pose{
			Pose: &gimock.PoseResult{
				RequestTime: timestamppb.Now(),
				PosX:        1.223,
				PosY:        2.223,
				PosZ:        3.223,
				QuatX:       4.223,
				QuatY:       5.223,
				QuatZ:       6.223,
				QuatW:       7.223,
				SlamState:   gimock.SlamState_TRACKING_GOOD,
			},
		},
	}, nil
}

func delay(t string) error {
	tint, err := strconv.Atoi(t)
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(tint) * time.Millisecond)
	return nil
}

func main() {
	port := os.Getenv("PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// TLSを使用しないという制約の中で、クライアント側の実装で、ベーシック認証のみを残すことができない。
	// そのため、サーバ側からもクライアント認証(basic認証)をはずしている。
	// https://github.com/machinemapplatform/mmpf-monolithic/issues/19
	s := grpc.NewServer()
	gimock.RegisterMmpfMonolithicServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func setMetadata(request, metadata map[string]string) map[string]string {
	for k, v := range request {
		metadata["metadata."+k] = v
	}
	return metadata
}

func setImages(images []*gimock.Image, metadata map[string]string) map[string]string {
	for k, v := range images {
		i := strconv.Itoa(k)
		metadata["images["+i+"].lens_placement"] = v.LensPlacement.Enum().String()
		metadata["images["+i+"].byte(length)"] = strconv.Itoa(len(v.GetByte()))
	}
	return metadata
}
