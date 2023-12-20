package types

func NewEventStream(eventName, channelID string) *EventStream {
	return &EventStream{
		EventName: eventName,
		ChannelId: channelID,
	}
}

func (EventStream) Validate() error {
	return nil
}

func NewInterchainEvent(name string, attributes []InterchainEventAttribute) *InterchainEvent {
	return &InterchainEvent{
		Name:       name,
		Attributes: attributes,
	}
}

func (InterchainEvent) Validate() error {
	return nil
}

func NewInterchainEventAttribute(name, value string) *InterchainEventAttribute {
	return &InterchainEventAttribute{
		Name:  name,
		Value: value,
	}
}

func (InterchainEventAttribute) Validate() error {
	return nil
}
