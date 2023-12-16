package types

import (
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
)

// DefaultGenesis creates and returns the default interchain query GenesisState
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Port:       PortID,
		Params:     DefaultParams(),
		Registered: []EventListener{},
		Listeners:  []EventListener{},
	}
}

// NewGenesisState creates a returns a new GenesisState instance
func NewGenesisState(hostPort string, params Params, registered, listeners []EventListener) *GenesisState {
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

	// TODO
	// For every registered listener
	// We must verify that each channel exists
	// We must verify that each event name is valid

	// TODO
	// For every registered listener
	// We must verify that each channel exists
	// We must verify that each event name is valid

	return gs.Params.Validate()
}
