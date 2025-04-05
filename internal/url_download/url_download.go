package url_download

import (
	"errors"
	"fmt"
	"github.com/samyarsadat/MC-Client-Prep/internal/common"
	"github.com/samyarsadat/MC-Client-Prep/internal/fs_helpers"
	"github.com/samyarsadat/MC-Client-Prep/internal/logger"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFromUrl(url, hash string, destination string, overwrite bool, progressChan chan<- DlProgress) (DlResult, error) {
	defer close(progressChan)

	logger.Log.Debug("URL download requested", "url", url, "hash", hash, "destination", destination)
	retRes := DlResult{filesize: -1}

	// We don't want to implicitly create the destination directory if it doesn't exist.
	destExists, err := fs_helpers.DirectoryExists(destination)
	if err != nil {
		return retRes, fmt.Errorf("fs error: %w", err)
	}
	if !destExists {
		return retRes, common.ErrDirNotExist
	}

	resp, err := http.Get(url)
	if err != nil {
		return retRes, fmt.Errorf("http request error: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return retRes, errors.New("http request failed with status code: " + resp.Status)
	}

	filename, err := getFileName(resp)
	if err != nil {
		return retRes, fmt.Errorf("error getting filename: %w", err)
	}

	fileDestination := filepath.Join(destination, filename)
	retRes.filename = filename
	retRes.filepath = fileDestination

	fileExists, err := fs_helpers.FileExists(fileDestination)
	if err != nil {
		return retRes, fmt.Errorf("error checking if file exists: %w", err)
	}
	if fileExists {
		if !overwrite {
			return retRes, common.ErrFileExistsNoOverwrite
		}

		logger.Log.Debug("file already exists, overwriting it", "file", fileDestination)

		err = os.Remove(fileDestination)
		if err != nil {
			return retRes, fmt.Errorf("error removing existing file: %w", err)
		}
	}

	file, err := os.Create(fileDestination)
	if err != nil {
		return retRes, fmt.Errorf("error creating file: %w", err)
	}

	// Content-Length can be -1 if the length is not known, so we must check for that.
	contentLength := resp.ContentLength
	retRes.filesize = contentLength

	if contentLength > 0 {
		progWriter := &DlProgressWriter{
			total:    uint64(contentLength),
			progChan: progressChan,
			filename: filename,
			filepath: fileDestination,
		}
		writer := io.MultiWriter(file, progWriter)
		_, err = io.Copy(writer, resp.Body)
	} else {
		_, err = io.Copy(file, resp.Body)
	}
	if err != nil {
		return retRes, fmt.Errorf("error copying response body to file: %w", err)
	}

	err = file.Close()
	if err != nil {
		_ = os.Remove(fileDestination) // At this point, it might as well fail.
		return retRes, fmt.Errorf("error closing file: %w", err)
	}

	if hash != "" {
		err = VerifyFileHash(fileDestination, hash)
		if err != nil {
			return retRes, fmt.Errorf("error verifying file hash: %w", err)
		}
	}

	logger.Log.Debug("successfully downloaded file", "file", file.Name(), "url", url)
	return retRes, nil
}

func getFileName(resp *http.Response) (string, error) {
	contentDisposition := resp.Header.Get("Content-Disposition")

	if contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			return "", fmt.Errorf("error parsing Content-Disposition header: %w", err)
		}

		filename := params["filename"]
		if filename != "" {
			return filename, nil
		}
	}

	return "", errors.New("filename not found in Content-Disposition header")
}
