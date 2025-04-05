package testutils

import (
	"fmt"
	"github.com/samyarsadat/MC-Client-Prep/internal/fs_helpers"
	"os"
	"path/filepath"
)

var tempDirSuffix = "goTestDir"

func GetCreateTempDir() (string, error) {
	destination := filepath.Join(os.TempDir(), tempDirSuffix)
	destExists, err := fs_helpers.DirectoryExists(destination)
	if err != nil {
		return "", fmt.Errorf("error checking if temp dir exists: %w", err)
	}
	if !destExists {
		err := os.Mkdir(destination, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("error creating temp dir: %w", err)
		}
		return destination, nil
	}

	err = fs_helpers.HasDirPerms(destination)
	if err != nil {
		return "", fmt.Errorf("temp dir permission error: %w", err)
	}

	return destination, nil
}

func DeleteTempDir() error {
	destination := filepath.Join(os.TempDir(), tempDirSuffix)
	err := os.RemoveAll(destination)
	return err
}
