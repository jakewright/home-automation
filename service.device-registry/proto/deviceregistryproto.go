package deviceregistryproto

const (
	// InterpolationContinuous describes an integer property that is updated
	// smoothly and continuously e.g. volume, as opposed to channel number.
	InterpolationContinuous = "continuous"
)

// DeviceHeader contains metadata about a device
type DeviceHeader struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Type           string                 `json:"type"` // e.g. huelight
	Kind           string                 `json:"kind"` // e.g. lamp
	Attributes     map[string]interface{} `json:"attributes"`
	RoomID         string                 `json:"room_id"`
	Room           *Room                  `json:"room,omitempty"`
	ControllerName string                 `json:"controller_name"`
	//StateProviders string
}

// Room is a physical room in the building
type Room struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Devices []*DeviceHeader `json:"devices,omitempty"`
}
