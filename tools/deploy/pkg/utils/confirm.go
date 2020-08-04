package utils

import (
	"fmt"

	"github.com/logrusorgru/aurora"

	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/tools/deploy/pkg/git"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

// Deployment describes a service deployment
type Deployment struct {
	ServiceName string
	ServicePath string
	TargetName  string
	TargetHost  string

	// CurrentRevision is the long git hash of the currently deployed revision.
	// This should be the empty string if the service hasn't been deployed yet.
	CurrentRevision string

	// NewRevision is the long git hash of the revision that is to be deployed
	NewRevision string
}

// ConfirmDeployment prompts the user to confirm deployment of a service
func ConfirmDeployment(d *Deployment) (bool, error) {
	newRevisionShort, err := git.ShortHash(d.NewRevision)
	if err != nil {
		return false, oops.WithMessage(err, "failed to get short hash of new revision")
	}

	var currentRevisionShort string
	if d.CurrentRevision != "" {
		currentRevisionShort, err = git.ShortHash(d.CurrentRevision)
		if err != nil {
			return false, oops.WithMessage(err, "failed to get short hash of current revision")
		}
	}

	fmt.Println() // blank line
	fmt.Printf("Service  %s\n", aurora.Index(105, d.ServiceName))
	fmt.Printf("Target   %s %s\n", aurora.Index(105, d.TargetName), aurora.Gray(16, d.TargetHost))
	if d.CurrentRevision == "" {
		fmt.Printf("Revision %s %s\n", aurora.Index(105, newRevisionShort), aurora.Gray(16, "(not currently deployed)"))
	} else {
		fmt.Printf("Revision %s\n", aurora.Sprintf(aurora.Index(105, "%s...%s"), currentRevisionShort, newRevisionShort))
	}
	fmt.Println() // blank line

	commits, err := git.Log(d.CurrentRevision, d.NewRevision, "./"+d.ServicePath)
	if err != nil {
		return false, oops.WithMessage(err, "failed to get commit list")
	}

	if len(commits) == 0 {
		fmt.Printf(aurora.Sprintf(aurora.Index(178, "No commits in this range for %s\n"), d.ServiceName))
	}

	for _, commit := range commits {
		fmt.Printf("%s %s\n", aurora.Index(178, commit.ShortHash), commit.TitleLine)
	}

	if !output.Confirm(true, "Continue?") {
		return false, nil
	}

	return true, nil
}
