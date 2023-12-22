package types

import (
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
)

// NewEventStream creates a new event stream.
func NewEventStream(eventName, channelID string) *EventStream {
	return &EventStream{
		EventName: eventName,
		ChannelId: channelID,
	}
}

// Validate checks the validity of the data on the event stream.
func (es EventStream) Validate() error {
	if err := host.ChannelIdentifierValidator(es.ChannelId); err != nil {
		return err
	}

	// TODO Check the event name ?

	return nil
}

// NewEventStream creates a new interchain event.
func NewInterchainEvent(name string, attributes []InterchainEventAttribute) *InterchainEvent {
	return &InterchainEvent{
		Name:       name,
		Attributes: attributes,
	}
}

// Validate checks the validity of the data on the interchain event.
func (event InterchainEvent) Validate() error {
	// TODO Check the name ?

	// Check each of the attributes
	for _, attr := range event.Attributes {
		if err := attr.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// NewInterchainEventAttribute creates a new interchain event attribute.
func NewInterchainEventAttribute(name, value string) *InterchainEventAttribute {
	return &InterchainEventAttribute{
		Name:  name,
		Value: value,
	}
}

// Validate checks the validity of the data on the interchain event attribute.
func (attr InterchainEventAttribute) Validate() error {
	// TODO Check the name
	// TODO Check the attribute
	return nil
}
