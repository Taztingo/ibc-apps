package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-apps/modules/ice/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
)

// GetListeners returns the streams listening for events from this chain.
func (k Keeper) GetListeners(ctx sdk.Context) []types.EventStream {
	events := []types.EventStream{}

	k.IterateUpstreamEvents(ctx, func(eventStream types.EventStream) (stop bool, err error) {
		events = append(events, eventStream)
		return
	})

	return events
}

// GetRegisteredEvents returns the streams that this chain is listening to events from.
func (k Keeper) GetRegisteredEvents(ctx sdk.Context) []types.EventStream {
	events := []types.EventStream{}

	k.IterateDownstreamEvents(ctx, func(eventStream types.EventStream) (stop bool, err error) {
		events = append(events, eventStream)
		return
	})

	return events
}

// RegisterDownstreamEvent adds a stream to listen to events from.
func (k Keeper) RegisterDownstreamEvent(ctx sdk.Context, event types.EventStream, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) error {
	if err := event.Validate(); err != nil {
		ctx.Logger().Error("failed to validate register downstream event: " + err.Error())
		return err
	}

	if k.HasDownstreamEvent(ctx, event.EventName, event.ChannelId) {
		ctx.Logger().Info(fmt.Sprintf("downstream event %s is already registered for channel %s", event.EventName, event.ChannelId))
		return types.ErrDownstreamEventFound
	}

	// TODO Check if channel exists

	k.SetDownstreamEvent(ctx, event)

	// Should not always be done
	_, err := k.SendRegisterEventPacket(ctx, event, "ice-listener", timeoutHeight, timeoutTimestamp)
	return err
}

// UnregisterDownstreamEvent removes a stream to listen to events from.
func (k Keeper) UnregisterDownstreamEvent(ctx sdk.Context, event types.EventStream, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) error {
	if err := event.Validate(); err != nil {
		ctx.Logger().Error("failed to validate unregister downstream event: " + err.Error())
		return err
	}

	if !k.HasDownstreamEvent(ctx, event.EventName, event.ChannelId) {
		ctx.Logger().Info(fmt.Sprintf("downstream event %s is not registered for channel %s", event.EventName, event.ChannelId))
		return types.ErrDownstreamEventNotFound
	}

	// TODO Check if channel exists

	k.RemoveDownstreamEvent(ctx, event.EventName, event.ChannelId)

	_, err := k.SendUnregisterEventPacket(ctx, event, "ice-listener", timeoutHeight, timeoutTimestamp)
	return err
}

// RegisterUpstreamEvent adds a stream to broadcast events to.
func (k Keeper) RegisterUpstreamEvent(ctx sdk.Context, event types.EventStream) error {
	if err := event.Validate(); err != nil {
		ctx.Logger().Error("failed to validate register upstream event: " + err.Error())
		return err
	}

	if k.HasUpstreamEvent(ctx, event.EventName, event.ChannelId) {
		ctx.Logger().Info(fmt.Sprintf("upstream event %s is already registered for channel %s", event.EventName, event.ChannelId))
		return types.ErrUpstreamEventFound
	}

	// TODO Check if channel exists

	k.SetUpstreamEvent(ctx, event)

	return nil
}

// UnregisterUpstreamEvent removes a stream to broadcast events to.
func (k Keeper) UnregisterUpstreamEvent(ctx sdk.Context, event types.EventStream) error {
	if err := event.Validate(); err != nil {
		ctx.Logger().Error("failed to validate unregister upstream event: " + err.Error())
		return err
	}

	if !k.HasUpstreamEvent(ctx, event.EventName, event.ChannelId) {
		ctx.Logger().Info(fmt.Sprintf("upstream event %s is not registered for channel %s", event.EventName, event.ChannelId))
		return types.ErrUpstreamEventNotFound
	}

	k.RemoveUpstreamEvent(ctx, event.EventName, event.ChannelId)

	return nil
}

// SetDownstreamEvent adds a downstream for an event to the store.
func (k Keeper) SetDownstreamEvent(ctx sdk.Context, event types.EventStream) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDownstreamEventKey(event.EventName, event.ChannelId)
	bz := k.cdc.MustMarshal(&event)
	store.Set(key, bz)
}

// HasDownstreamEvent checks if the store has a downstream for the event name.
func (k Keeper) HasDownstreamEvent(ctx sdk.Context, eventName, channelID string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDownstreamEventKey(eventName, channelID)
	return store.Has(key)
}

// RemoveDownstreamEvent removes a downstream for an event from the store.
func (k Keeper) RemoveDownstreamEvent(ctx sdk.Context, eventName, channelID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDownstreamEventKey(eventName, channelID)
	if store.Has(key) {
		store.Delete(key)
	}
}

// RemoveDownstreamEvent gets a downstream for an event from the store.
func (k Keeper) GetDownstreamEvent(ctx sdk.Context, eventName, channelID string) (stream types.EventStream, err error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDownstreamEventKey(eventName, channelID)
	bz := store.Get(key)
	if len(bz) == 0 {
		return stream, types.ErrDownstreamEventNotFound
	}
	err = k.cdc.Unmarshal(bz, &stream)
	return stream, err
}

// IterateDownstreamEvents iterates through each of the downstreams the chain is listening to.
func (k Keeper) IterateDownstreamEvents(ctx sdk.Context, handle func(stream types.EventStream) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DownstreamEventPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		record := types.EventStream{}
		if err := k.cdc.Unmarshal(iterator.Value(), &record); err != nil {
			return err
		}
		stop, err := handle(record)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

// SetUpstreamEvent adds an upstream for an event to the store.
func (k Keeper) SetUpstreamEvent(ctx sdk.Context, event types.EventStream) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUpstreamEventKey(event.EventName, event.ChannelId)
	bz := k.cdc.MustMarshal(&event)
	store.Set(key, bz)
}

// HasUpstreamEvent checks if the store has an upstream for the event name.
func (k Keeper) HasUpstreamEvent(ctx sdk.Context, eventName, channelID string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUpstreamEventKey(eventName, channelID)
	return store.Has(key)
}

// GetUpstreamEvent gets an upstream for an event from the store.
func (k Keeper) GetUpstreamEvent(ctx sdk.Context, eventName, channelID string) (stream types.EventStream, err error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUpstreamEventKey(eventName, channelID)
	bz := store.Get(key)
	if len(bz) == 0 {
		return stream, types.ErrUpstreamEventNotFound
	}
	err = k.cdc.Unmarshal(bz, &stream)
	return stream, err
}

// RemoveUpstreamEvent removes an upstream for an event from the store.
func (k Keeper) RemoveUpstreamEvent(ctx sdk.Context, eventName, channelID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUpstreamEventKey(eventName, channelID)
	if store.Has(key) {
		store.Delete(key)
	}
}

// IterateUpstreamEvents iterates through each of the upstreams the chain is broadcasting to.
func (k Keeper) IterateUpstreamEvents(ctx sdk.Context, handle func(eventStream types.EventStream) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UpstreamEventPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		record := types.EventStream{}
		if err := k.cdc.Unmarshal(iterator.Value(), &record); err != nil {
			return err
		}
		stop, err := handle(record)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}
