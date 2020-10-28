package routes

import (
	"os"
	"testing"

	"github.com/jakewright/home-automation/libraries/go/bootstrap"
)

func TestMain(m *testing.M) {
	bootstrap.SetupTest()
	os.Exit(m.Run())
}
