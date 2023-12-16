package keeper

import (
	"github.com/cosmos/ibc-apps/modules/ice/interchain-query-demo/x/interquery/types"
)

var _ types.QueryServer = Keeper{}
