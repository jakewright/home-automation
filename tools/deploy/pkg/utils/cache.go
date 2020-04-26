package utils

import (
	"os"
	"path"
)

// CacheDir returns the directory that can be used
// as a temporary working directory by the tool.
func CacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	return path.Join(dir, "home-automation")
}
