package types

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrUnknownDataType         = sdkerrors.Register(ModuleName, 1, "unknown data type")
	ErrInvalidChannelFlow      = sdkerrors.Register(ModuleName, 2, "invalid message sent to channel end")
	ErrInvalidHostPort         = sdkerrors.Register(ModuleName, 3, "invalid host port")
	ErrHostDisabled            = sdkerrors.Register(ModuleName, 4, "host is disabled")
	ErrInvalidVersion          = sdkerrors.Register(ModuleName, 5, "invalid version")
	ErrDownstreamEventNotFound = sdkerrors.Register(ModuleName, 6, "downstream event does not exist")
	ErrUpstreamEventNotFound   = sdkerrors.Register(ModuleName, 7, "upstream event does not exist")
	ErrUpstreamEventFound      = sdkerrors.Register(ModuleName, 8, "upstream event already exists")
	ErrDownstreamEventFound    = sdkerrors.Register(ModuleName, 9, "downstream event already exists")
)
