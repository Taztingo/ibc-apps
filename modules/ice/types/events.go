package types

// ICQ Interchain Event events
const (
	EventTypePacketError     = "icq_packet_error"
	EventTypeInterchainEvent = "interchain_event"

	AttributeKeyAckError      = "error"
	AttributeKeyHostChannelID = "host_channel_id"
	AttributeKeyEventName     = "event_name"
	AttributeKeyChainID       = "chain_id"
)
