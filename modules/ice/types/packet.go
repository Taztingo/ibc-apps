package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateBasic performs basic validation of the interchain query packet data.
func (packet InterchainRegisterPacket) ValidateBasic() error {
	return nil
}

// GetBytes returns the JSON marshalled interchain query packet data.
func (packet InterchainRegisterPacket) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&packet))
}

// ValidateBasic performs basic validation of the interchain query packet data.
func (packet InterchainUnregisterPacket) ValidateBasic() error {
	return nil
}

// GetBytes returns the JSON marshalled interchain query packet data.
func (packet InterchainUnregisterPacket) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&packet))
}

// ValidateBasic performs basic validation of the interchain query packet data.
func (packet InterchainEventPacket) ValidateBasic() error {
	return nil
}

// GetBytes returns the JSON marshalled interchain query packet data.
func (packet InterchainEventPacket) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&packet))
}
