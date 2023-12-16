package keeper_test

import (
	"github.com/cosmos/ibc-apps/modules/ice/testing/simapp"
	"github.com/cosmos/ibc-apps/modules/ice/types"
)

func (suite *KeeperTestSuite) TestInitGenesis() {
	suite.SetupTest()

	genesisState := types.GenesisState{
		HostPort: TestPort,
		Params: types.Params{
			HostEnabled: false,
		},
	}

	simapp.GetSimApp(suite.chainA).ICQKeeper.InitGenesis(suite.chainA.GetContext(), genesisState)

	port := simapp.GetSimApp(suite.chainA).ICQKeeper.GetPort(suite.chainA.GetContext())
	suite.Require().Equal(TestPort, port)

	expParams := types.NewParams(
		false,
	)
	params := simapp.GetSimApp(suite.chainA).ICQKeeper.GetParams(suite.chainA.GetContext())
	suite.Require().Equal(expParams, params)
}

func (suite *KeeperTestSuite) TestExportGenesis() {
	suite.SetupTest()

	genesisState := simapp.GetSimApp(suite.chainA).ICQKeeper.ExportGenesis(suite.chainA.GetContext())

	suite.Require().Equal(types.PortID, genesisState.GetHostPort())

	expParams := types.DefaultParams()
	suite.Require().Equal(expParams, genesisState.GetParams())
}
