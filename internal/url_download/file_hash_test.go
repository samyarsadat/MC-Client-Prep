package url_download

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/samyarsadat/MC-Client-Prep/internal/testutils"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestHashFileSHA1(t *testing.T) {
	err := testutils.DeleteTempDir()
	assert.NoError(t, err)
	tempDir, err := testutils.GetCreateTempDir()
	assert.NoError(t, err)
	defer func() { _ = testutils.DeleteTempDir() }()

	fileName := "testfile-" + testutils.GetRandomBase64(10)
	filePath := filepath.Join(tempDir, fileName)

	mockedContent := []byte(testutils.GetRandomBase64(1000)) // 1KB
	file, err := os.Create(filePath)
	assert.NoError(t, err)
	_, err = file.Write(mockedContent)
	assert.NoError(t, err)
	err = file.Close()
	assert.NoError(t, err)

	hasher := sha1.New()
	hasher.Write(mockedContent)
	act_hash := hex.EncodeToString(hasher.Sum(nil))

	hash, err := HashFileSHA1(filePath)
	assert.NoError(t, err)
	assert.Equal(t, act_hash, hash)
}
