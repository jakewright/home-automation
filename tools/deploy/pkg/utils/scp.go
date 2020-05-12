package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

// SCP initiates an scp transfer from src to dst
func SCP(src, username, host, dst string) error {
	var args []string

	if fi, err := os.Stat(src); err != nil {
		return oops.WithMessage(err, "failed to stat src")
	} else if fi.IsDir() {
		args = append(args, "-r")
	}

	args = append(args, src, fmt.Sprintf("%s@%s:%s", username, host, dst))

	output.Debug("scp %s", strings.Join(args, " "))

	if err := exe.Command("scp", args...).Run().Err; err != nil {
		return oops.WithMessage(err, "failed to scp file")
	}

	return nil
}
