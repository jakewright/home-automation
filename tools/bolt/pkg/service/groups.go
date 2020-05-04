package service

// Group represents a predefined group of services
type Group struct {
	Name     string
	Services []string
}

// Groups is a list of groups that can be specified on the command line
var Groups = []*Group{
	{
		Name:     "core",
		Services: []string{"service.api-gateway", "service.config", "service.device-registry", "redis", "mysql"},
	},
}
