package domain

// Actor represents a device to perform an action on
type Actor struct {
	Identifier     string `json:"identifier"`
	ControllerName string `json:"controller_name"`
}
