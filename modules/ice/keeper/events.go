package keeper

import (
	"github.com/cosmos/ibc-apps/modules/ice/types"
	icetypes "github.com/cosmos/ibc-apps/modules/ice/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

// EmitWriteErrorAcknowledgementEvent emits an event signalling an error acknowledgement and including the error details
func EmitWriteErrorAcknowledgementEvent(ctx sdk.Context, packet exported.PacketI, err error) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			icetypes.EventTypePacketError,
			sdk.NewAttribute(sdk.AttributeKeyModule, icetypes.ModuleName),
			sdk.NewAttribute(icetypes.AttributeKeyAckError, err.Error()),
			sdk.NewAttribute(icetypes.AttributeKeyHostChannelID, packet.GetDestChannel()),
		),
	)
}

func EmitInterchainEvent(ctx sdk.Context, ice types.InterchainEvent, chainID string) {
	attributes := make([]sdk.Attribute, len(ice.Attributes)+2)
	attributes = append(attributes, sdk.NewAttribute(icetypes.AttributeKeyEventName, ice.Name))
	attributes = append(attributes, sdk.NewAttribute(icetypes.AttributeKeyChainID, chainID))
	for _, attribute := range ice.Attributes {
		attributes = append(attributes, sdk.NewAttribute(attribute.Name, attribute.Value))
	}
	event := sdk.NewEvent(icetypes.EventTypeInterchainEvent, attributes...)
	ctx.EventManager().EmitEvent(event)
}
