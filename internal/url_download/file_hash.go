package url_download

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func HashFileSHA1(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %s", err)
	}
	defer func() { _ = file.Close() }()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("error hashing file: %s", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func VerifyFileHash(filePath string, hash string) error {
	fileHash, err := HashFileSHA1(filePath)
	if err != nil {
		return err
	}
	if fileHash != hash {
		return fmt.Errorf("file hash does not match expected hash")
	}

	return nil
}
