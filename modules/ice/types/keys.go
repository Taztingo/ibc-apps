package types

const (
	// ModuleName defines the interchain query module name
	ModuleName = "interchainevents"

	// PortID is the default port id that the interchain events module binds to
	PortID = "icehost"

	// Version defines the current version for interchain events
	Version = "ice-1"

	// StoreKey is the store key string for interchain query
	StoreKey = ModuleName

	// RouterKey is the message route for interchain query
	RouterKey = ModuleName

	// QuerierRoute is the querier route for interchain query
	QuerierRoute = ModuleName
)

var (
	// ParamsKey defines the key to store the params in store
	ParamsKey = []byte{0x00}
	// PortKey defines the key to store the port ID in store
	PortKey = []byte{0x01}
	// DownstreamEventPrefix defines the prefix all downstream events are stored with in the store
	DownstreamEventPrefix = []byte{0x02}
	// UpstreamEventPrefix defines the prefix all upstream events are stored with in the store
	UpstreamEventPrefix = []byte{0x03}
)

func GetUpstreamEventKey(eventName string) []byte {
	key := UpstreamEventPrefix
	key = append(key, []byte(eventName)...)
	return key
}

func GetDownstreamEventKey(eventName string) []byte {
	key := DownstreamEventPrefix
	key = append(key, []byte(eventName)...)
	return key
}
