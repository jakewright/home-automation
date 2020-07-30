package domain

// Action is a single change to make to a device
type Action struct {
	ScheduleID int
	Property   string
	Value      string
	Type       string
}
