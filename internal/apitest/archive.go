package apitest

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// CreateZipFileFromDirectory creates a zip file from a directory, the zip file will be generated into a
// temporary directory
func CreateZipFileFromDirectory(t *testing.T, source string) string {
	temp, err := os.MkdirTemp("", strings.ReplaceAll(source, "/", "-")+"-*")
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err := os.RemoveAll(temp)
		if err != nil {
			t.Error(err)
		}
	}()

	// Set the target file path within the temporary directory
	targetFile := filepath.Join(temp, "archive.zip")

	// Create a new zip file to write to
	zipFile, err := os.Create(targetFile)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err := zipFile.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	// Create a new zip archive
	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		err := zipWriter.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	// Walk through the source directory and add each file to the zip archive
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and hidden files
		if info.IsDir() || filepath.Base(path)[0] == '.' {
			return nil
		}

		// Add the file to the zip archive
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			err := file.Close()
			if err != nil {
				t.Error(err)
			}
		}()

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Error(err)
	}

	return targetFile
}
