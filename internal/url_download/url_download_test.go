package url_download

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/samyarsadat/MC-Client-Prep/internal/logger"
	"github.com/samyarsadat/MC-Client-Prep/internal/testutils"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetFileName(t *testing.T) {
	resp := httpmock.NewStringResponse(200, "")
	fileName := "testfile-" + testutils.GetRandomBase64(10)
	resp.Header.Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	name, err := getFileName(resp)
	assert.NoError(t, err)
	assert.Equal(t, fileName, name)
}

func TestDownloadFromUrl(t *testing.T) {
	loggerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger.Init(loggerOpts)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := "https://testurl.com"
	mockedContent := []byte(testutils.GetRandomBase64(1000)) // 1KB
	progressChan := make(chan DlProgress)

	hasher := sha1.New()
	hasher.Write(mockedContent)
	hash := hex.EncodeToString(hasher.Sum(nil))

	err := testutils.DeleteTempDir()
	assert.NoError(t, err)
	tempDir, err := testutils.GetCreateTempDir()
	assert.NoError(t, err)
	defer func() { _ = testutils.DeleteTempDir() }()

	fileName := "testfile-" + testutils.GetRandomBase64(10)
	filePath := filepath.Join(tempDir, fileName)

	// 500 bytes/sec speed limit
	thrReader := &testutils.ThrottledReader{
		Reader:    bytes.NewReader(mockedContent),
		ChunkSize: 10,
		Delay:     20 * time.Millisecond,
	}

	contentLength := int64(len(mockedContent))
	httpmock.RegisterResponder("GET", url, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewBytesResponse(200, mockedContent)
		resp.Body = io.NopCloser(thrReader)
		resp.Header.Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		resp.ContentLength = contentLength
		return resp, nil
	})

	go func() {
		result, err := DownloadFromUrl(url, hash, tempDir, true, progressChan)
		assert.NoError(t, err)
		assert.Equal(t, fileName, result.filename)
		assert.Equal(t, filePath, result.filepath)
		assert.Equal(t, contentLength, result.filesize)
	}()

	var oldProgPct uint8 = 0
	for progress := range progressChan {
		fmt.Printf("\rDownloading %s to %s... %d%%\n", progress.filename, progress.filepath, progress.progPct)
		assert.Greater(t, progress.progPct, oldProgPct)
		assert.LessOrEqual(t, progress.progPct, uint8(100))
		assert.Equal(t, uint8(5), progress.progPct-oldProgPct)
		oldProgPct = progress.progPct
	}
	assert.Equal(t, uint8(100), oldProgPct)

	data, _ := os.ReadFile(filePath)
	assert.Equal(t, mockedContent, data)
}
