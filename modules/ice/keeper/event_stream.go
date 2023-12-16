package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-apps/modules/ice/types"
)

// GetListeners returns the module's listeners.
func (k Keeper) GetListeners(_ sdk.Context) []types.EventStream {
	return []types.EventStream{}
}

// GetRegisteredEvents returns the module's registered events.
func (k Keeper) GetRegisteredEvents(_ sdk.Context) []types.EventStream {
	return []types.EventStream{}
}

func (k Keeper) RegisterEvent(ctx sdk.Context, event, channel string) {

}

func (k Keeper) UnregisterEvent(ctx sdk.Context, event, channel string) {

}

func (k Keeper) PublishEvent(ctx sdk.Context, event string, attributes ...string) {

}
