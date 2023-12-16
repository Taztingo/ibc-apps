package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-apps/modules/ice/types"
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
func (k Keeper) RegisterDownstreamEvent(ctx sdk.Context, event types.EventStream) error {
	if err := event.Validate(); err != nil {
		return err
	}

	if k.HasDownstreamEvent(ctx, event.EventName) {
		return types.ErrDownstreamEventFound
	}

	// Check if channel exists

	k.SetDownstreamEvent(ctx, event)

	return nil
}

// UnregisterDownstreamEvent removes a stream to listen to events from.
func (k Keeper) UnregisterDownstreamEvent(ctx sdk.Context, event types.EventStream) error {
	if err := event.Validate(); err != nil {
		return err
	}

	if !k.HasDownstreamEvent(ctx, event.EventName) {
		return types.ErrDownstreamEventNotFound
	}

	k.RemoveDownstreamEvent(ctx, event.EventName)

	return nil
}

// RegisterUpstreamEvent adds a stream to broadcast events to.
func (k Keeper) RegisterUpstreamEvent(ctx sdk.Context, event types.EventStream) error {
	if err := event.Validate(); err != nil {
		return err
	}

	if k.HasDownstreamEvent(ctx, event.EventName) {
		return types.ErrUpstreamEventFound
	}

	// Check if channel exists

	k.SetUpstreamEvent(ctx, event)

	return nil
}

// UnregisterUpstreamEvent removes a stream to broadcast events to.
func (k Keeper) UnregisterUpstreamEvent(ctx sdk.Context, event types.EventStream) error {
	if err := event.Validate(); err != nil {
		return err
	}

	if !k.HasUpstreamEvent(ctx, event.EventName) {
		return types.ErrUpstreamEventNotFound
	}

	k.RemoveUpstreamEvent(ctx, event.EventName)

	return nil
}

// SetDownstreamEvent adds a downstream for an event to the store.
func (k Keeper) SetDownstreamEvent(ctx sdk.Context, event types.EventStream) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDownstreamEventKey(event.EventName)
	bz := k.cdc.MustMarshal(&event)
	store.Set(key, bz)
}

// HasDownstreamEvent checks if the store has a downstream for the event name.
func (k Keeper) HasDownstreamEvent(ctx sdk.Context, eventName string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDownstreamEventKey(eventName)
	return store.Has(key)
}

// RemoveDownstreamEvent removes a downstream for an event from the store.
func (k Keeper) RemoveDownstreamEvent(ctx sdk.Context, eventName string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDownstreamEventKey(eventName)
	if store.Has(key) {
		store.Delete(key)
	}
}

// RemoveDownstreamEvent gets a downstream for an event from the store.
func (k Keeper) GetDownstreamEvent(ctx sdk.Context, eventName string) (stream types.EventStream, err error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDownstreamEventKey(eventName)
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
	key := types.GetUpstreamEventKey(event.EventName)
	bz := k.cdc.MustMarshal(&event)
	store.Set(key, bz)
}

// HasUpstreamEvent checks if the store has an upstream for the event name.
func (k Keeper) HasUpstreamEvent(ctx sdk.Context, eventName string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUpstreamEventKey(eventName)
	return store.Has(key)
}

// GetUpstreamEvent gets an upstream for an event from the store.
func (k Keeper) GetUpstreamEvent(ctx sdk.Context, eventName string) (stream types.EventStream, err error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUpstreamEventKey(eventName)
	bz := store.Get(key)
	if len(bz) == 0 {
		return stream, types.ErrUpstreamEventNotFound
	}
	err = k.cdc.Unmarshal(bz, &stream)
	return stream, err
}

// RemoveUpstreamEvent removes an upstream for an event from the store.
func (k Keeper) RemoveUpstreamEvent(ctx sdk.Context, eventName string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUpstreamEventKey(eventName)
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
