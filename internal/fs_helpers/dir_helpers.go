package fs_helpers

import (
	"errors"
	"github.com/samyarsadat/MC-Client-Prep/internal/common"
	"io/fs"
	"os"
	"path/filepath"
)

func DirectoryExists(path string) (bool, error) {
	_, err := os.Open(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func canWriteToDir(path string) error {
	testFile := filepath.Join(path, ".tmp_perm_test")
	file, err := os.Create(testFile)
	if err != nil {
		return err
	}
	_ = file.Close()
	_ = os.Remove(testFile)
	return nil
}

func canReadDir(path string) error {
	_, err := os.ReadDir(path)
	return err
}

func HasDirPerms(path string) error {
	exists, err := DirectoryExists(path)
	if err != nil {
		return err
	}
	if !exists {
		return common.ErrDirNotExist
	}

	err = canWriteToDir(path)
	if err != nil {
		return err
	}

	return canReadDir(path)
}
