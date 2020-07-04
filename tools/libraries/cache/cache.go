package cache

import (
	"fmt"
	"os"
	"path"
)

var dir string

// Init must be called before Dir() can be used
func Init(tool string) error {
	osCacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("failed to get user cache dir: %w", err)
	}

	dir = path.Join(osCacheDir, "home-automation", tool)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create %s: %w", dir, err)
	}

	return nil
}

// Dir returns the directory that can be used
// as a temporary working directory by the tool.
func Dir() string {
	if dir == "" {
		panic(fmt.Errorf("function Dir() called before cache.Init()"))
	}

	return dir
}
