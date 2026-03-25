package cache

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zemelkajakub/lskv/internal/config"
)

const cacheSubDir = "cache"

// GetCacheDir returns the cache directory path
func GetCacheDir() (string, error) {
	appDirPath, err := config.GetAppDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(appDirPath, cacheSubDir)

	return cacheDir, nil
}

// EnsureCacheDir ensures that cache directory exists and returns its path
func EnsureCacheDir() error {
	cacheDirPath, err := GetCacheDir()
	if err != nil {
		return err
	}

	// Ensure .lskv directory exists
	err = os.MkdirAll(cacheDirPath, 0o700)
	if err != nil {
		return fmt.Errorf("error creating '%s' directory: %w", cacheDirPath, err)
	}
	return nil
}
