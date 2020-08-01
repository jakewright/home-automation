package bootstrap

import "github.com/jakewright/home-automation/libraries/go/distsync"

// SetupTest should be called in a TestMain() function to
// setup the various global state that code relies on.
func SetupTest() {
	distsync.DefaultLocksmith = distsync.NewLocalLocksmith()
}
