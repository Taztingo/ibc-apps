package keeper_test

import (
	"github.com/cosmos/ibc-apps/modules/ice/testing/simapp"
	"github.com/cosmos/ibc-apps/modules/ice/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestQueryParams() {
	ctx := sdk.WrapSDKContext(suite.chainA.GetContext())
	expParams := types.DefaultParams()
	res, _ := simapp.GetSimApp(suite.chainA).ICQKeeper.Params(ctx, &types.QueryParamsRequest{})
	suite.Require().Equal(&expParams, res.Params)
}
