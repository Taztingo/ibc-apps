package keeper

import (
	"github.com/cosmos/ibc-apps/modules/ice/types"

	"cosmossdk.io/errors"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
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
	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		// Nothing needs to be done
	case *channeltypes.Acknowledgement_Error:
		k.AttemptRollbackRegisterEventPacket(ctx, packet)
		k.AttemptRollbackUnregisterEventPacket(ctx, packet)
	}

	return nil
}

func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet) error {
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

	// Do Logic
	// What do we want to do when we receive an event packet?
	// We either use the callback middleware or we handle callbacks ourselves
	// If we don't do it ourselves then Callback Middleware needs information about the packet.
	// We could store last processed packet type

	return types.NewInterchainPacketAck(types.InterchainPacketAck_EVENT).GetBytes(), nil
}

func (k Keeper) AttemptRegisterEventPacket(ctx sdk.Context, packet channeltypes.Packet) ([]byte, error) {
	var data types.InterchainRegisterPacket

	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		// UnmarshalJSON errors are indeterminate and therefore are not wrapped and included in failed acks
		return nil, errors.Wrapf(types.ErrUnknownDataType, "cannot unmarshal ICQ packet data")
	}

	event := types.EventStream{
		EventName: data.Event,
		ChannelId: packet.DestinationChannel,
	}
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

	k.RemoveDownstreamEvent(ctx, data.Event)
	return true
}

func (k Keeper) AttemptRollbackUnregisterEventPacket(ctx sdk.Context, packet channeltypes.Packet) bool {
	var data types.InterchainRegisterPacket
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return false
	}

	event := types.EventStream{
		EventName: data.Event,
		ChannelId: packet.DestinationChannel,
	}
	k.SetDownstreamEvent(ctx, event)
	return true
}

func (k Keeper) AttemptUnregisterEventPacket(ctx sdk.Context, packet channeltypes.Packet) ([]byte, error) {
	var data types.InterchainUnregisterPacket

	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		// UnmarshalJSON errors are indeterminate and therefore are not wrapped and included in failed acks
		return nil, errors.Wrapf(types.ErrUnknownDataType, "cannot unmarshal ICQ packet data")
	}

	event := types.EventStream{
		EventName: data.Event,
		ChannelId: packet.DestinationChannel,
	}
	err := k.UnregisterUpstreamEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return types.NewInterchainPacketAck(types.InterchainPacketAck_UNREGISTER).GetBytes(), nil
}

func (k Keeper) BroadcastEvent(ctx sdk.Context, event types.InterchainEvent) error {
	if err := event.Validate(); err != nil {
		return err
	}

	listeners := k.GetListeners(ctx)
	for _, listener := range listeners {
		_, err := k.SendEventPacket(ctx, event, listener.ChannelId, types.PortID, clienttypes.ZeroHeight(), 0)
		if err != nil {
			// Log here
		}
	}

	return nil
}

func (k Keeper) SendEventPacket(ctx sdk.Context, event types.InterchainEvent, sourceChannel, sourcePort string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) (uint64, error) {

	/*channel, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return 0, errorsmod.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}*/

	// destinationPort := channel.GetCounterparty().GetPortID()
	// destinationChannel := channel.GetCounterparty().GetChannelID()

	// begin createOutgoingPacket logic
	// See spec for this logic: https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#packet-relay
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, errorsmod.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packetData := types.InterchainEventPacket{
		Event: event,
	}
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

	packetData := types.InterchainRegisterPacket{
		Event: event.EventName,
	}
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

	packetData := types.InterchainUnregisterPacket{
		Event: event.EventName,
	}
	sequence, err := k.ics4Wrapper.SendPacket(ctx, channelCap, sourcePort, event.ChannelId, timeoutHeight, timeoutTimestamp, packetData.GetBytes())
	if err != nil {
		return 0, err
	}

	return sequence, nil
}
