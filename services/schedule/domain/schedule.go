package domain

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Schedule wraps a set of rules and a set of actions
type Schedule struct {
	gorm.Model

	// ActorID is the ID of the controller to act upon
	ActorID string

	// Actions is the list of actions to perform
	Actions []Action

	// StartTime is the earliest time that the schedule can run.
	// N.b. it might not run at this time if the rules do not permit.
	StartTime time.Time

	// NextRun is a cache of the next run time
	NextRun time.Time

	// Count is the number of times the schedule should run.
	// A value of -1 will run the schedule ad infinitum.
	Count int

	// Until is the end date of the schedule
	Until time.Time

	// Rules define when this schedule should run
	//Rules []ScheduleRule
}
