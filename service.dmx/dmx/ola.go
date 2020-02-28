package dmx

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/jakewright/home-automation/libraries/go/errors"
)

// OLA sends DMX information via the ola_set_dmx program
type OLA struct {
	m sync.Mutex
}

// Set sets all of the DMX values for the given universe
func (o *OLA) Set(universe int, values [512]byte) error {
	o.m.Lock()
	defer o.m.Unlock()

	a := args(universe, values)

	cmd := exec.Command("ola_set_dmx", a...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return errors.WithMessage(err, "failed to run ola_set_dmx %s", strings.Join(a, " "))
	}

	// Assume anything written to stderr is a bad thing
	if stderr.Len() > 0 {
		return errors.InternalService("ola_set_dmx wrote to stderr: %s", stderr.String())
	}

	return nil
}

func args(universe int, values [512]byte) []string {
	// Iterate over the values backwards until we
	// find the first non-zero value. This should
	// make sure we're not actually sending 512
	// arguments to ola_set_dmx.
	var slice []byte
	for i := 511; i >= 0; i-- {
		if values[i] > 0 {
			break
		}
		slice = values[:i]
	}

	args := []string{"--universe", strconv.Itoa(universe), "--dmx"}
	for _, v := range slice {
		args = append(args, strconv.Itoa(int(v)))
	}

	return args
}
