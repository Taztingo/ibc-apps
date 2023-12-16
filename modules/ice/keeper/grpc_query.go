package keeper

import (
	"context"

	"github.com/cosmos/ibc-apps/modules/ice/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = Keeper{}

// Params implements the Query/Params gRPC method
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{
		Params: &params,
	}, nil
}

// Listeners implements types.QueryServer.
func (k Keeper) Listeners(c context.Context, _ *types.QueryListenersRequest) (*types.QueryListenersResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	listeners := k.GetListeners(ctx)
	return &types.QueryListenersResponse{
		Listeners: listeners,
	}, nil
}

// RegisteredEvents implements types.QueryServer.
func (k Keeper) RegisteredEvents(c context.Context, _ *types.QueryRegisteredEventsRequest) (*types.QueryRegisteredEventsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	events := k.GetRegisteredEvents(ctx)
	return &types.QueryRegisteredEventsResponse{
		Registered: events,
	}, nil
}
