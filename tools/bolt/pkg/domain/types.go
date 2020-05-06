package domain

// Service represents a service
type Service struct {
	Name  string
	Ports []*Port
}

// Port describes a port mapping
type Port struct {
	Host      string
	Container string
}
