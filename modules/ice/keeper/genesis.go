package keeper

import (
	"fmt"

	"github.com/cosmos/ibc-apps/modules/async-icq/v7/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the icq state and binds to PortID.
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	k.SetPort(ctx, state.Port)

	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	if !k.IsBound(ctx, state.Port) {
		// transfer module binds to the transfer port on InitChain
		// and claims the returned capability
		err := k.BindPort(ctx, state.Port)
		if err != nil {
			panic(fmt.Sprintf("could not claim port capability: %v", err))
		}
	}

	if err := k.SetParams(ctx, state.Params); err != nil {
		panic(fmt.Sprintf("could not set params: %v", err))
	}

	// We need to set the registered
	// We need to set the listeners
}

// ExportGenesis exports icq module's portID and denom trace info into its genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Port:   k.GetPort(ctx),
		Params: k.GetParams(ctx),
	}
}
