package toolutils

import (
	"fmt"
	"os"
	"path"
)

var dir string

// Init must be called before CacheDir() can be used
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

// CacheDir returns the directory that can be used
// as a temporary working directory by the tool.
func CacheDir() string {
	if dir == "" {
		panic(fmt.Errorf("CacheDir() called before toolutils.Init()"))
	}

	return dir
}
