package listing

import (
	pb "booking/internal/proto/listing"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RpcListingClient struct {
	client pb.ListingServiceClient
	//conn   *grpc.ClientConn
}

const layer = "listing infra"

func NewRpcListingClient(adr string) (*RpcListingClient, error) {
	conn, err := grpc.NewClient(adr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewListingServiceClient(conn)
	return &RpcListingClient{client: client}, nil

}
func (c *RpcListingClient) PrepareGetListingRequest(id string) *pb.ListingRequest {
	return &pb.ListingRequest{
		Id: id,
	}
}

func (c *RpcListingClient) GetListing(ctx context.Context, in *pb.ListingRequest, opts ...grpc.CallOption) (*pb.ListingResponse, error) {
	res, err := c.client.GetListing(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("[%s] get listing failed: %v", layer, err)
	}
	return res, nil
}

//func (c *RpcListingClient) Close() error {
//	return c.conn.Close()
//}
