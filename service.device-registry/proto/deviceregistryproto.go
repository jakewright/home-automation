package deviceregistryproto

import "encoding/json"

const (
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

type ListDeviceResponse []*DeviceHeader

type GetDeviceResponse *DeviceHeader

type ListRoomsResponse []*Room

type GetRoomResponse *Room

type Property interface {
	MarshalJSON() ([]byte, error)
}

type IntProperty struct {
	Value         int    `json:"value"`
	Min           *int   `json:"value,omitempty"`
	Max           *int   `json:"value,omitempty"`
	Interpolation string `json:",omitempty"`
}

func (p *IntProperty) MarshalJSON() ([]byte, error) {
	type Alias IntProperty
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "int",
		Alias: (*Alias)(p),
	})
}
