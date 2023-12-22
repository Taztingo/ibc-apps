package types

import (
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
)

// DefaultGenesis creates and returns the default interchain query GenesisState
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Port:       PortID,
		Params:     DefaultParams(),
		Registered: []EventStream{},
		Listeners:  []EventStream{},
	}
}

// NewGenesisState creates a returns a new GenesisState instance
func NewGenesisState(hostPort string, params Params, registered, listeners []EventStream) *GenesisState {
	return &GenesisState{
		Port:       hostPort,
		Params:     params,
		Registered: registered,
		Listeners:  listeners,
	}
}

// Validate performs basic validation of the GenesisState
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.Port); err != nil {
		return err
	}

	for _, registered := range gs.Registered {
		if err := registered.Validate(); err != nil {
			return err
		}
	}

	for _, listener := range gs.Listeners {
		if err := listener.Validate(); err != nil {
			return err
		}
	}

	return gs.Params.Validate()
}
