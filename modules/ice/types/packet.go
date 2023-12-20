package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewInterchainRegisterPacket creates a new InterchainRegisterPacket.
func NewInterchainRegisterPacket(event, memo string) *InterchainRegisterPacket {
	return &InterchainRegisterPacket{
		Event: event,
		Memo:  memo,
	}
}

// ValidateBasic performs basic validation of the interchain register packet data.
func (packet InterchainRegisterPacket) ValidateBasic() error {
	return nil
}

// GetBytes returns the JSON marshalled interchain register packet data.
func (packet InterchainRegisterPacket) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&packet))
}

// NewInterchainUnregisterPacket creates a new InterchainUnregisterPacket.
func NewInterchainUnregisterPacket(event, memo string) *InterchainUnregisterPacket {
	return &InterchainUnregisterPacket{
		Event: event,
		Memo:  memo,
	}
}

// ValidateBasic performs basic validation of the interchain unregister packet data.
func (packet InterchainUnregisterPacket) ValidateBasic() error {
	return nil
}

// GetBytes returns the JSON marshalled interchain unregister packet data.
func (packet InterchainUnregisterPacket) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&packet))
}

// NewInterchainEventPacket creates a new InterchainEventPacket.
func NewInterchainEventPacket(event InterchainEvent, memo string) *InterchainEventPacket {
	return &InterchainEventPacket{
		Event: event,
		Memo:  memo,
	}
}

// ValidateBasic performs basic validation of the interchain event packet data.
func (packet InterchainEventPacket) ValidateBasic() error {
	return nil
}

// GetBytes returns the JSON marshalled interchain event packet data.
func (packet InterchainEventPacket) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&packet))
}

// NewInterchainPacketAck creates a new InterchainPacketAck.
func NewInterchainPacketAck(packetType InterchainPacketAck_Type) InterchainPacketAck {
	return InterchainPacketAck{
		Type: packetType,
	}
}

// GetBytes returns the JSON marshalled interchain ack packet data.
func (ack InterchainPacketAck) GetBytes() []byte {
	return ModuleCdc.MustMarshalJSON(&ack)
}

// ValidateBasic performs basic validation of the interchain ack packet data.
func (ack InterchainPacketAck) ValidateBasic() error {
	return nil
}
