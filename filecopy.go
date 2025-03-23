package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var ErrFileExistsAndIsDifferent = errors.New("destination file exists and is different")

func hashFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func CopyFile(src, dstDir string, overwrite bool) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file '%s' (%w)", src, err)
	}
	defer srcFile.Close()
	dstPath := filepath.Join(dstDir, filepath.Base(src))
	if _, err := os.Stat(dstPath); err == nil {
		srcHash, err1 := hashFile(src)
		dstHash, err2 := hashFile(dstPath)
		// If both hashes match, the files are identical
		if err1 == nil && err2 == nil && string(srcHash) == string(dstHash) {
			if dbg != nil {
				dbg.Info(fmt.Sprintf("CopyFile: Destination file is identical, skipping copy, source=%s, dest=%s", src, dstPath))
			}
			return nil
		} else {
			if !overwrite {
				if dbg != nil {
					dbg.Error(fmt.Sprintf("CopyFile: Destination file already exists and is different from source and overwriting is disabled, skipping copy, source=%s, dest=%s", src, dstPath))
				}
				return ErrFileExistsAndIsDifferent
			} else {
				if dbg != nil {
					dbg.Warn(fmt.Sprintf("CopyFile: Destination file already exists and is different from source and overwriting is enabled, will try to overwrite it, source=%s, dest=%s", src, dstPath))
				}
			}
		}
	}

	// Open the destination file properly, avoiding truncation issues.
	dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create destination file '%s' (%w)", dstPath, err)
	}
	defer func() {
		dstFile.Sync() // Ensure data is flushed before closing
		dstFile.Close()
	}()
	if bytesWritten, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file '%s' to '%s' (%w)", src, dstPath, err)
	} else {
		if dbg != nil {
			dbg.Info(fmt.Sprintf("CopyFile: File copied successfully, bytes written=%d, source=%s, dest=%s", bytesWritten, src, dstPath))
		}
	}
	return nil
}
