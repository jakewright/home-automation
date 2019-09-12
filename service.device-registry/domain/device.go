package domain

type Device struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Type           string                 `json:"type"` // e.g. huelight
	Kind           string                 `json:"kind"` /// e.g. lamp
	Attributes     map[string]interface{} `json:"attributes"`
	RoomID         string                 `json:"room_id"`
	Room           *Room                  `json:"room,omitempty"`
	ControllerName string                 `json:"controller_name"`
	//StateProviders string
}
