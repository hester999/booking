package listing

import (
	client "booking/internal/infra/rpc/listing"
	"booking/internal/models"
	"context"
)

type ListingGateway struct {
	client *client.RpcListingClient
}

func NewListingGateway(client *client.RpcListingClient) *ListingGateway {

	return &ListingGateway{
		client: client,
	}
}

func (g *ListingGateway) GetListing(ctx context.Context, id string) (*models.Listing, error) {
	req := g.client.PrepareGetListingRequest(id)

	result, err := g.client.GetListing(ctx, req)
	if err != nil {
		return nil, err
	}
	return &models.Listing{
		ID:                     result.Id,
		MinimumBookingDuration: result.MinimumBookingDuration,
		Price:                  result.Price,
	}, nil

}
