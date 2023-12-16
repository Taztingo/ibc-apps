package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/cosmos/ibc-apps/modules/ice/interchain-query-demo/testutil/keeper"
	"github.com/cosmos/ibc-apps/modules/ice/interchain-query-demo/x/interquery/keeper"
	"github.com/cosmos/ibc-apps/modules/ice/interchain-query-demo/x/interquery/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.InterqueryKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
