package keeper

import (
	"fmt"

	"github.com/cosmos/ibc-apps/modules/ice/types"

	"cosmossdk.io/errors"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	tendermintclient "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

// OnRecvPacket handles a given interchain queries packet on a destination host chain.
// If the transactFion is successfully executed, the transaction response bytes will be returned.
func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet) ([]byte, error) {

	if ack, err := k.AttemptRecvEventPacket(ctx, packet); err == nil {
		return ack, err
	}
	if ack, err := k.AttemptRegisterEventPacket(ctx, packet); err == nil {
		return ack, err
	}
	if ack, err := k.AttemptUnregisterEventPacket(ctx, packet); err == nil {
		return ack, err
	}

	return nil, errors.Wrapf(types.ErrUnknownDataType, "cannot unmarshal ICE packet data")
}

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, ack channeltypes.Acknowledgement) error {
	switch res := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		// Nothing needs to be done
	case *channeltypes.Acknowledgement_Error:
		ctx.Logger().Error(fmt.Sprintf("received ack error for packet %v: %s", packet, res.Error))
		k.AttemptRollbackRegisterEventPacket(ctx, packet)
		k.AttemptRollbackUnregisterEventPacket(ctx, packet)
	}

	return nil
}

func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet) error {
	ctx.Logger().Info(fmt.Sprintf("received timeout packet for packet: %v", packet))
	k.AttemptRollbackRegisterEventPacket(ctx, packet)
	k.AttemptRollbackUnregisterEventPacket(ctx, packet)
	return nil
}

func (k Keeper) AttemptRecvEventPacket(ctx sdk.Context, packet channeltypes.Packet) ([]byte, error) {
	var data types.InterchainEventPacket

	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		// UnmarshalJSON errors are indeterminate and therefore are not wrapped and included in failed acks
		return nil, errors.Wrapf(types.ErrUnknownDataType, "cannot unmarshal ICE packet data")
	}

	// We emit the event
	chainID := k.GetChainID(ctx, packet.SourcePort, packet.SourceChannel)
	EmitInterchainEvent(ctx, data.Event, chainID)

	// We leave it up to the Callback module to handle the rest

	return types.NewInterchainPacketAck(types.InterchainPacketAck_EVENT).GetBytes(), nil
}

func (k Keeper) AttemptRegisterEventPacket(ctx sdk.Context, packet channeltypes.Packet) ([]byte, error) {
	var data types.InterchainRegisterPacket

	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		// UnmarshalJSON errors are indeterminate and therefore are not wrapped and included in failed acks
		return nil, errors.Wrapf(types.ErrUnknownDataType, "cannot unmarshal ICQ packet data")
	}

	event := *types.NewEventStream(data.Event, packet.DestinationChannel)
	err := k.RegisterUpstreamEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return types.NewInterchainPacketAck(types.InterchainPacketAck_REGISTER).GetBytes(), nil
}

func (k Keeper) AttemptRollbackRegisterEventPacket(ctx sdk.Context, packet channeltypes.Packet) bool {
	var data types.InterchainRegisterPacket
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return false
	}

	ctx.Logger().Info(fmt.Sprintf("rollback: removing downstream event %s on channel %s for packet %v", data.Event, packet.DestinationChannel, packet))

	k.RemoveDownstreamEvent(ctx, data.Event, packet.DestinationChannel)
	return true
}

func (k Keeper) AttemptRollbackUnregisterEventPacket(ctx sdk.Context, packet channeltypes.Packet) bool {
	var data types.InterchainUnregisterPacket
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return false
	}

	ctx.Logger().Info(fmt.Sprintf("rollback: adding downstream event %s on channel %s for packet %v", data.Event, packet.DestinationChannel, packet))

	event := *types.NewEventStream(data.Event, packet.DestinationChannel)
	k.SetDownstreamEvent(ctx, event)
	return true
}

func (k Keeper) AttemptUnregisterEventPacket(ctx sdk.Context, packet channeltypes.Packet) ([]byte, error) {
	var data types.InterchainUnregisterPacket

	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		// UnmarshalJSON errors are indeterminate and therefore are not wrapped and included in failed acks
		return nil, errors.Wrapf(types.ErrUnknownDataType, "cannot unmarshal ICE packet data")
	}

	ctx.Logger().Info(fmt.Sprintf("unregistering event %s on channel %s", data.Event, packet.DestinationChannel))

	event := *types.NewEventStream(data.Event, packet.DestinationChannel)
	err := k.UnregisterUpstreamEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return types.NewInterchainPacketAck(types.InterchainPacketAck_UNREGISTER).GetBytes(), nil
}

func (k Keeper) BroadcastEvent(ctx sdk.Context, event types.InterchainEvent) error {
	if err := event.Validate(); err != nil {
		ctx.Logger().Error("failed to validate interchain event: " + err.Error())
		return err
	}

	listeners := k.GetListeners(ctx)
	ctx.Logger().Info(fmt.Sprintf("broadcasting event: %v to listeners: %v", event, listeners))
	for _, listener := range listeners {
		_, err := k.SendEventPacket(ctx, event, listener.ChannelId, types.PortID, clienttypes.ZeroHeight(), 0)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("failed to send packet to channel %v on port %v. error: %v", listener.ChannelId, types.PortID, err))
		}
	}

	return nil
}

func (k Keeper) SendEventPacket(ctx sdk.Context, event types.InterchainEvent, sourceChannel, sourcePort string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) (uint64, error) {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, errorsmod.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packetData := types.NewInterchainEventPacket(event, "")
	sequence, err := k.ics4Wrapper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packetData.GetBytes())
	if err != nil {
		return 0, err
	}

	return sequence, nil
}

func (k Keeper) SendRegisterEventPacket(ctx sdk.Context, event types.EventStream, sourcePort string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) (uint64, error) {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, event.ChannelId))
	if !ok {
		return 0, errorsmod.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packetData := types.NewInterchainRegisterPacket(event.EventName, "")
	sequence, err := k.ics4Wrapper.SendPacket(ctx, channelCap, sourcePort, event.ChannelId, timeoutHeight, timeoutTimestamp, packetData.GetBytes())
	if err != nil {
		return 0, err
	}

	return sequence, nil
}

func (k Keeper) SendUnregisterEventPacket(ctx sdk.Context, event types.EventStream, sourcePort string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) (uint64, error) {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, event.ChannelId))
	if !ok {
		return 0, errorsmod.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packetData := types.NewInterchainUnregisterPacket(event.EventName, "")
	sequence, err := k.ics4Wrapper.SendPacket(ctx, channelCap, sourcePort, event.ChannelId, timeoutHeight, timeoutTimestamp, packetData.GetBytes())
	if err != nil {
		return 0, err
	}

	return sequence, nil
}

func (k Keeper) GetChainID(ctx sdk.Context, ibcPort, ibcChannel string) string {
	chainID := "unknown"
	_, clientState, err := k.channelKeeper.GetChannelClientState(ctx, ibcPort, ibcChannel)
	if err != nil {
		return chainID
	}

	tmClientState, ok := clientState.(*tendermintclient.ClientState)
	if ok {
		return tmClientState.ChainId
	}
	return chainID
}
