package domain

// Room is a physical room in the building
type Room struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Devices []*Device `json:"devices,omitempty"`
}
