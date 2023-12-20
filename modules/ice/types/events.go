package types

// ICQ Interchain Event events
const (
	EventTypePacketError     = "ice_packet_error"
	EventTypeInterchainEvent = "interchain_event"

	AttributeKeyAckError      = "error"
	AttributeKeyHostChannelID = "host_channel_id"
	AttributeKeyEventName     = "event_name"
	AttributeKeyChainID       = "chain_id"
)
