package build

import (
	"os"
	"path"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/util"
)

const (
	repository      = "git@github.com:jakewright/home-automation.git"
	mirrorDirectory = "mirror"
)

func Checkout(reference string) error {
	if err := cacheDir(); err != nil {
		return errors.Wrap(err, "failed to create cache directory")
	}

	// If the clone directory does not exist
	if _, err := os.Stat(mirrorDirectory); err != nil && os.IsNotExist(err) {
		// Clone the repository
		if err := util.Exec("git", "clone", repository, mirrorDirectory); err != nil {
			return errors.Wrap(err, "failed to clone repository")
		}
	}

	// Change the working directory to the git clone
	if err := os.Chdir(mirrorDirectory); err != nil {
		return errors.Wrap(err, "failed to change working directory to mirror directory")
	}

	// Checkout the commit or branch
	if err := util.Exec("git", "checkout", reference); err != nil {
		return errors.Wrap(err, "failed to checkout reference '%s'", reference)
	}

	return nil
}

func cacheDir() error {
	dir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	dir = path.Join(dir, "home-automation")

	// Use MkdirAll instead of Mkdir because it
	// doesn't error if the directory already exists
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	if err := os.Chdir(dir); err != nil {
		return err
	}

	return nil
}
