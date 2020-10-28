package routes

import "github.com/jakewright/home-automation/libraries/go/database"

// Controller handles requests
type Controller struct {
	Database database.Database
}
