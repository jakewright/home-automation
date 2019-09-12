package domain

type Room struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Devices []*Device `json:"devices,omitempty"`
}
