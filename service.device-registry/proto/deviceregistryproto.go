package deviceregistryproto

// Device is a thing that can be controlled by the system
type Device struct {
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
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Devices []*Device `json:"devices,omitempty"`
}

type ListDeviceResponse []*Device

type GetDeviceResponse *Device

type ListRoomsResponse []*Room

type GetRoomResponse *Room
