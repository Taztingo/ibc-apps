package keeper

import (
	"github.com/cosmos/ibc-apps/modules/ice/types"
	icqtypes "github.com/cosmos/ibc-apps/modules/ice/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

// EmitWriteErrorAcknowledgementEvent emits an event signalling an error acknowledgement and including the error details
func EmitWriteErrorAcknowledgementEvent(ctx sdk.Context, packet exported.PacketI, err error) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			icqtypes.EventTypePacketError,
			sdk.NewAttribute(sdk.AttributeKeyModule, icqtypes.ModuleName),
			sdk.NewAttribute(icqtypes.AttributeKeyAckError, err.Error()),
			sdk.NewAttribute(icqtypes.AttributeKeyHostChannelID, packet.GetDestChannel()),
		),
	)
}

func EmitInterchainEvent(ctx sdk.Context, ice types.InterchainEvent, chainID string) {
	attributes := make([]sdk.Attribute, len(ice.Attributes)+1)
	attributes = append(attributes, sdk.NewAttribute("chain-id", chainID))
	for _, attribute := range ice.Attributes {
		attributes = append(attributes, sdk.NewAttribute(attribute.Name, attribute.Value))
	}
	event := sdk.NewEvent(ice.Name, attributes...)
	ctx.EventManager().EmitEvent(event)
}
