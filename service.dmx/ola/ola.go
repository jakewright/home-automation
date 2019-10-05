package ola

import (
	"bytes"
	"os/exec"
	"strconv"

	"github.com/jakewright/home-automation/libraries/go/errors"
)

// SetDMX sets all of the DMX values for the given universe
func SetDMX(universe int, values [512]byte) error {
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

	cmd := exec.Command("ola_set_dmx", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Assume anything written to stderr is a bad thing
	if stderr.Len() > 0 {
		return errors.InternalService("ola_set_dmx wrote to stderr: %s", stderr.String())
	}

	return nil
}
