package types_test

import (
	"testing"

	"github.com/cosmos/ibc-apps/modules/async-icq/v7/types"
	"github.com/stretchr/testify/suite"

	ibctesting "github.com/cosmos/ibc-go/v7/testing"
)

type TypesTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
}

func (suite *TypesTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)

	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

func (suite *TypesTestSuite) TestValidateGenesisState() {
	var genesisState types.GenesisState

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {},
			true,
		},
		{
			"failed to validate - empty value",
			func() {
				genesisState = types.GenesisState{}
			},
			false,
		},
		{
			"failed to validate - invalid host port",
			func() {
				genesisState = *types.NewGenesisState("p", types.DefaultParams(), []types.EventListener{}, []types.EventListener{})
			},
			false,
		},
		{
			"failed to validate - invalid event name on registered",
			func() {
				genesisState = *types.NewGenesisState("port", types.NewParams(true), []types.EventListener{}, []types.EventListener{})
			},
			false,
		},
		{
			"failed to validate - invalid channel on registered",
			func() {
				genesisState = *types.NewGenesisState("port", types.NewParams(true), []types.EventListener{}, []types.EventListener{})
			},
			false,
		},
		{
			"failed to validate - invalid event name on listeners",
			func() {
				genesisState = *types.NewGenesisState("port", types.NewParams(true), []types.EventListener{}, []types.EventListener{})
			},
			false,
		},
		{
			"failed to validate - invalid channel on listeners",
			func() {
				genesisState = *types.NewGenesisState("port", types.NewParams(true), []types.EventListener{}, []types.EventListener{})
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			genesisState = *types.DefaultGenesis()

			tc.malleate() // malleate mutates test data

			err := genesisState.Validate()

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}
